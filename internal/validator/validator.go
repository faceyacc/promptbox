package validator

import (
	"net/mail"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	// Holds error messages for form field.
	FieldErrors    map[string]string
	NonFieldErrors []string
}

// Valid() checks if any validation and non validations errors exist.
// If no errors exist Valid() returns true.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map
// if an entry does not exist already for a given field.
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldErrors(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField() adds an error message to the FieldErrors map only
// if a validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, mesage string) {
	if !ok {
		v.AddFieldError(key, mesage)
	}
}

// NotBlank() returns true if an string is not empty.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if an string is no longer than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue() returns true if value T is in a list of specific
// permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// MinChars returns true if value contains a least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches returns true if a value matches a given regular expression
func Matches(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil
}
