package services

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Response struct {
	Status   string `json:"status"` // "passed" or "failed"
	Error    string `json:"error,omitempty"`
	Res      string `json:"res,omitempty"`
	TestCase string `json:"testcase"`
}

func cleanUpTempDir(dir string) {
	os.Remove(dir)
}

func createDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return cli, nil
}

func removeDockerContainer(client *client.Client, contID string) error {
	err := client.ContainerRemove(context.Background(), contID, container.RemoveOptions{
		Force: true,
	})
	return err
}

func removeDockerImage(client *client.Client, imageName string) error {
	_, err := client.ImageRemove(context.Background(), imageName, image.RemoveOptions{
		Force: true,
	})
	return err
}

func getCodeFileName(language string) (string, error) {
	var fileName string

	switch language {
	case "python":
		fileName = "main.py"
	case "c++":
		fileName = "main.cpp"
	case "java":
		fileName = "main.java"
	case "javascript":
		fileName = "main.js"
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
	case "c++":
		Dockerfile = `
			FROM gcc:latest
			WORKDIR /app
			COPY . /app
			CMD ["sh", "-c", "g++ -o main main.cpp && ./main"]
			`
	case "java":
		Dockerfile = `
			FROM openjdk:17-slim
			WORKDIR /app
			COPY . /app
			CMD ["sh", "-c", "javac main.java && java Main"]
			`
	case "javascript":
		Dockerfile = `
			FROM node:20
			WORKDIR /app
			COPY . /app
			CMD ["node", "main.js"]
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
	case "c++":
		imageName = fmt.Sprintf("%s-%s-image", baseName, "cpp")
	case "java":
		imageName = fmt.Sprintf("%s-%s-image", baseName, "java")
	case "javascript":
		imageName = fmt.Sprintf("%s-%s-image", baseName, "javascript")
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

func runContainer(client *client.Client, dir, dockerImageName, containerName, port string, inputEnv []string) (string, error) {
	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func ExecuteCodeInDocker(language, code string) ([]Response, error) {
	// create a temporary dir to contain file and Dockerfile inside it
	dir, err := os.MkdirTemp("", "code")
	if err != nil {
		return []Response{}, err
	}
	defer cleanUpTempDir(dir)

	// get code file name based on language
	fileName, err := getCodeFileName(language)
	if err != nil {
		return []Response{}, err
	}

	// write code in this file
	filePath := fmt.Sprintf("%s/%s", dir, fileName)
	err = os.WriteFile(filePath, []byte(code), 0744)
	if err != nil {
		return []Response{}, err
	}

	// get dockerfile content based on language
	dockerfileContent, err := getDockerfileContent(language)
	if err != nil {
		return []Response{}, err
	}

	// write dockerfile in this dir
	dockerfilePath := fmt.Sprintf("%s/Dockerfile", dir)
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0744)
	if err != nil {
		return []Response{}, err
	}

	// create docker client
	client, err := createDockerClient()
	if err != nil {
		return []Response{}, err
	}

	// build Image from dockerfile
	dockerImageName := getImageName(language)
	tags := []string{dockerImageName}
	err = buildImageFromDockerfile(client, tags, dockerfilePath)
	if err != nil {
		return []Response{}, err
	}

	// run the container from image
	containerName := strings.Replace(dockerImageName, "image", "container", 1)
	containerOutput, err := runContainer(client, dir, dockerImageName, containerName, "8080", []string{})
	if err != nil {
		return []Response{}, err
	}
	log.Printf("Container output: %s", containerOutput)

	err = removeDockerImage(client, dockerImageName)
	if err != nil {
		return []Response{}, err
	}

	return []Response{}, nil
}