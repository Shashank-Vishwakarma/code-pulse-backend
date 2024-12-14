package utils

import (
	"github.com/go-playground/validator/v10"
)

func ValidateRequest(body interface{}) error {
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return err
	}
	return nil
}
