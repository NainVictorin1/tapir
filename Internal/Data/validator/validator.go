package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Improved email regex
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Validator struct {
	Errors map[string]string
}

// NewValidator initializes a Validator
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// ValidData checks if there are no validation errors
func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

// AddError records an error for a given field
func (v *Validator) AddError(field string, message string) {
	v.Errors[field] = message
}

// ValidateField checks a condition and adds an error if false
func (v *Validator) ValidateField(valid bool, field, message string) {
	if !valid {
		v.AddError(field, message)
	}
}

// NotBlank checks if a string is not empty or only spaces
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinLength checks if a string meets the minimum length requirement
func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxLength checks if a string does not exceed the maximum length
func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// IsValidEmail validates an email address format
func IsValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}
