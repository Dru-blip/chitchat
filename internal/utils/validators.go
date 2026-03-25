package utils

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (cv *Validator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
