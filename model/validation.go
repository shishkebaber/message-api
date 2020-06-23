package model

import (
	"fmt"
	"github.com/go-playground/validator"
	"time"
)

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Error: '%s' field validation failed on the '%s' tag",
		v.Field(),
		v.Tag(),
	)
}

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("timestamp", customTimestampValidation)
	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)
	var resultErrs ValidationErrors
	if errs != nil {
		errs := errs.(validator.ValidationErrors)
		if len(errs) == 0 {
			return nil
		}

		for _, err := range errs {
			e := ValidationError{err.(validator.FieldError)}
			resultErrs = append(resultErrs, e)
		}
	}

	return resultErrs
}

// Custom validation for timestamp field
func customTimestampValidation(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02 15:04:05.000", fl.Field().String())
	if err != nil {
		return false
	}

	return true
}
