package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Validation is a struct for validation
type Validation struct {
	validate *validator.Validate
}

// NewValidator creates a new Validation type
func NewValidator() *Validation {

	v := validator.New()
	// Register custom validation for sku field
	v.RegisterValidation("sku", validateSKU)

	return &Validation{validate: v}
}

// Validate used for struct validation
func (v *Validation) Validate(i interface{}) ValidationErrors {

	// Custom validate struct to wrap validation errors that expose field error
	vErr := v.validate.Struct(i)
	if vErr != nil {
		var returnErrs []ValidationError
		for _, err := range vErr.(validator.ValidationErrors) {
			ve := ValidationError{err.(validator.FieldError)}
			returnErrs = append(returnErrs, ve)
		}
		return returnErrs
	}

	return nil

	// vErr := v.validate.Struct(i).(validator.ValidationErrors)

	// if len(vErr) == 0 {
	// 	return nil
	// }

	// var returnErrs []ValidationError
	// for _, err := range vErr {
	// 	ve := ValidationError{err.(validator.FieldError)}
	// 	returnErrs = append(returnErrs, ve)
	// }
	// return returnErrs
}

// ValidationError wraps the validators FieldError so we do not
// expose this to out code
type ValidationError struct {
	validator.FieldError
}

// ValidationErrors is an error collection of ValidationError
type ValidationErrors []ValidationError

func (ve ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		ve.Namespace(),
		ve.Field(),
		ve.Tag(),
	)
}

// Errors converts the slice into a string slice
func (ves ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range ves {
		errs = append(errs, err.Error())
	}

	return errs
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku format : abc-abc-xy
	re := regexp.MustCompile(`[a-z]{3}-[0-9]{3}-[a-z]{3}`)
	match := re.FindAllString(fl.Field().String(), -1)

	if len(match) == 1 {
		return true
	}

	return false
}
