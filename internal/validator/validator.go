package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldsErrors []string
	FieldsErrors    map[string]string
}

// Valid() returns true if the FieldsErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldsErrors) == 0 && len(v.NonFieldsErrors) == 0
}

// AddFieldError() adds an error message to the FieldsErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already
	// initialized.
	if v.FieldsErrors == nil {
		v.FieldsErrors = make(map[string]string)
	}
	if _, exists := v.FieldsErrors[key]; !exists {
		v.FieldsErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldsErrors = append(v.NonFieldsErrors, message)
}

// CheckField() adds an error message to the FieldsErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// MinChars() returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// EqualValue() returns true if a value is equals expected value
// values.
func EqualValue[T comparable](value T, expectedValue T) bool {
	return value == expectedValue
}

// PermittedValue() returns true if a value is in a list of specific permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches() returns true if a value matches a specified pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsEmail(value string) bool {
	return Matches(value, EmailRX)
}
