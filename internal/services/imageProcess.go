package services

import (
	"context"
	"os"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func generateUid() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		logrus.Errorf("Failed to generate UUID: %v", err)
		return ""
	}

	return uuid.String()
}

func UploadImageToCloudinary(image string) (string, error) {
	// Remove image from local
	defer func() {
		os.Remove(image)
	}()

	cld, err := cloudinary.NewFromParams(config.Config.CLOUDINARY_CLOUD_NAME, config.Config.CLOUDINARY_API_KEY, config.Config.CLOUDINARY_API_SECRET)
	if err != nil {
		return "", err
	}

	result, err := cld.Upload.Upload(context.Background(), image, uploader.UploadParams{PublicID: "blog_image" + "-" + image + "-" + generateUid()})
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
