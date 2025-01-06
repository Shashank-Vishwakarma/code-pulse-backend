package utils

import (
	"fmt"
	"mime/multipart"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateRequest(body interface{}) error {
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return err
	}
	return nil
}

func ValidateImageFile(image *multipart.FileHeader) error {
	// check if the file is an image
	acceptableMimeTypes := []string{"jpeg", "png", "jpg", "webp"}
	mimeType := strings.Split(image.Filename, ".")[len(strings.Split(image.Filename, "."))-1]
	if !slices.Contains(acceptableMimeTypes, mimeType) {
		return fmt.Errorf("file is not an image")
	}

	// check for file size
	if image.Size > 5*1024*1024 { // 5MB
		return fmt.Errorf("file is too large. Please try to uploada file smaller than 5MB")
	}

	return nil
}
