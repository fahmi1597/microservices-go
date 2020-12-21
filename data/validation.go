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

	validate := validator.New()
	// Register custom validation for sku field
	validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku format : abc-123-xy
	fs := fl.Field().String()
	re := regexp.MustCompile(`[a-z]{3}-[0-9]{3}-[a-z]{2}`)
	m := re.FindAllString(fs, -1)

	if len(m) != 1 {
		return false
	}

	return true
}

// Validate is a struct validator
func (v *Validation) Validate(p interface{}) ValidationErrors {

	// vErr := v.validate.Struct(p)

	// if vErr != nil {
	// 	return vErr
	// }

	// return nil

	// Custom error message
	vErr := v.validate.Struct(p)

	if vErr != nil {
		var returnErrs []ValidationError
		for _, err := range vErr.(validator.ValidationErrors) {
			ve := ValidationError{err.(validator.FieldError)}
			returnErrs = append(returnErrs, ve)
		}
		return returnErrs
	}

	return nil

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
