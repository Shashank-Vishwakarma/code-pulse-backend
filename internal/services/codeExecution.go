package services

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Response struct {
	Input string `json:"input,omitempty"`
	Output      string `json:"output,omitempty"`
	Expected string `json:"expected,omitempty"`
	Result   bool   `json:"result,omitempty"`
}

func cleanUpTempDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to delete temp dir %s: %w", dir, err)
	}

	return nil
}

func removeDockerContainer(client *client.Client, contID string) error {
	// Remove the main container
	err := client.ContainerRemove(context.Background(), contID, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete main container %s: %w", contID, err)
	}

	// Remove containers which are left in "created" state
	containers, err := client.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("status", "created"),
		),
	})
	if err != nil {
		return fmt.Errorf("failed to list created containers: %w", err)
	}

	// Loop and remove these containers
	for _, cont := range containers {
		err := client.ContainerRemove(context.Background(), cont.ID, container.RemoveOptions{
			Force: true,
		})
		if err != nil {
			return fmt.Errorf("failed to delete container %s: %w", cont.ID, err)
		}
	}

	return nil
}

func removeDockerImage(client *client.Client, imageName string) error {
	// Remove the main image
	_, err := client.ImageRemove(context.Background(), imageName, image.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete main image %s: %w", imageName, err)
	}

	images, err := client.ImageList(context.Background(), image.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("dangling", "true"),
		),
	})
	if err != nil {
		return fmt.Errorf("failed to list dangling images: %w", err)
	}

	// Loop and remove each dangling image
	for _, img := range images {
		_, err := client.ImageRemove(context.Background(), img.ID, image.RemoveOptions{
			Force: true,
		})
		if err != nil {
			return fmt.Errorf("failed to delete dangling image %s: %w", img.ID, err)
		}
	}

	return nil
}

func getCodeFileName(language string) (string, error) {
	var fileName string

	switch language {
	case "python":
		fileName = "main.py"
	case "javascript":
		fileName = "main.js"
	case "go":
		fileName = "main.go"
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	return fileName, nil
}

func getDockerfileContent(language string) (string, error) {
	var Dockerfile string

	switch language {
	case "python":
		Dockerfile = `
			FROM python:3.12-slim
			WORKDIR /app
			COPY . /app
			CMD ["python", "main.py"]
			`
	case "javascript":
		Dockerfile = `
			FROM node:20-slim
			WORKDIR /app
			COPY . /app
			CMD ["node", "main.js"]
			`
	case "go":
		Dockerfile = `
			FROM golang:alpine
			WORKDIR /app
			COPY . /app
			CMD ["go", "run", "main.go"]
			`
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	return Dockerfile, nil
}

func getImageName(language string) string {
	var imageName string

	baseName := "codepulse-code-execution"

	switch language {
	case "python":
		imageName = fmt.Sprintf("%s-%s-image", baseName, "python")
	case "javascript":
		imageName = fmt.Sprintf("%s-%s-image", baseName, "javascript")
	case "go":
		imageName = fmt.Sprintf("%s-%s-image", baseName, "go")
	default:
		imageName = ""
	}

	return imageName
}

func buildImageFromDockerfile(client *client.Client, tags []string, dockerfilePath string) error {
	ctx := context.Background()

	buffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buffer)
	defer tarWriter.Close()

	// read dockerfile
	dockerfileReader, err := os.Open(dockerfilePath)
	if err != nil {
		return err
	}
	dockerfileContent, err := io.ReadAll(dockerfileReader)
	if err != nil {
		return err
	}

	// write dockerfile content to tar file
	tarHeader := &tar.Header{
		Name: dockerfilePath,
		Size: int64(len(dockerfileContent)),
	}

	err = tarWriter.WriteHeader(tarHeader)
	if err != nil {
		return err
	}

	_, err = tarWriter.Write(dockerfileContent)
	if err != nil {
		return err
	}

	dockerTarFileReader := bytes.NewReader(buffer.Bytes())

	buildOptions := types.ImageBuildOptions{
		Context: dockerTarFileReader,
		Dockerfile: dockerfilePath,
		Tags: tags,
	}

	// Build the actual image
	imageBuildResponse, err := client.ImageBuild(ctx, dockerTarFileReader, buildOptions)
	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()

	// Read the STDOUT from the build process
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func runContainer(client *client.Client, dir, dockerImageName, containerName string) (string, error) {
	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	// create container
	cont, err := client.ContainerCreate(
		ctx, 
		&container.Config{
			Image: dockerImageName,
		},
		&container.HostConfig{
			// mount the volume [temp dir -> docker container] to run e.g. python main.py 
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: dir,
					Target: "/app",
				},
			},
			// limit resources
			Resources: container.Resources{
				Memory: 512 * 1024 * 1024,
				CPUQuota: 50000,
			},
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return "", err
	}
	
	defer removeDockerContainer(client, cont.ID)

	// start the container
	err = client.ContainerStart(ctx, cont.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	statusCh, errCh := client.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
		case err := <-errCh:
			return "", fmt.Errorf("error while waiting for container: %w", err)
		case <-ctx.Done():
			// Timeout occured, stop the container
			err = client.ContainerStop(ctx, cont.ID, container.StopOptions{})
			if err != nil {
				return "", fmt.Errorf("failed to stop container after timeout: %w", err)
			}
			return "", fmt.Errorf("execution timed out")
		case <-statusCh:
	}

	// get the logs from container
	out, err := client.ContainerLogs(ctx, cont.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}

	output, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func createDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli, nil
}

func ExecuteCodeInDocker(language, code string) (string, error) {
	// create a temporary dir to contain file and Dockerfile inside it
	dir, err := os.MkdirTemp("", "code")
	if err != nil {
		return "", err
	}
	defer cleanUpTempDir(dir)

	// get code file name based on language
	fileName, err := getCodeFileName(language)
	if err != nil {
		return "", err
	}

	// write code in this file
	filePath := fmt.Sprintf("%s/%s", dir, fileName)
	err = os.WriteFile(filePath, []byte(code), 0744)
	if err != nil {
		return "", err
	}

	// get dockerfile content based on language
	dockerfileContent, err := getDockerfileContent(language)
	if err != nil {
		return "", err
	}

	// write dockerfile in this dir
	dockerfilePath := fmt.Sprintf("%s/Dockerfile", dir)
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0744)
	if err != nil {
		return "", err
	}

	// create docker client
	client, err := createDockerClient()
	if err != nil {
		return "", err
	}

	// build Image from dockerfile
	dockerImageName := getImageName(language)
	tags := []string{dockerImageName}
	err = buildImageFromDockerfile(client, tags, dockerfilePath)
	if err != nil {
		return "", err
	}

	// run the container from image
	containerName := strings.Replace(dockerImageName, "image", "container", 1)
	containerOutput, err := runContainer(client, dir, dockerImageName, containerName)
	if err != nil {
		return "", err
	}

	output := regexp.MustCompile(`�\[`).ReplaceAllString(containerOutput, "[")
	output = regexp.MustCompile(`'`).ReplaceAllString(output, `"`)
	output = regexp.MustCompile(`\b(True|False)\b`).ReplaceAllStringFunc(output, func(match string) string {
		if match == "True" {
			return "true"
		}
		return "false"
	})
	output = regexp.MustCompile(`"output": (\d+)`).ReplaceAllString(output, `"output": "$1"`)

	err = removeDockerImage(client, dockerImageName)
	if err != nil {
		return "", err
	}

	return output, nil
}