package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsNotEmpty tests the IsNotEmpty validator.
func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"nil value", nil, errors.New("value is nil")},
		{"empty string", "", errors.New("value is empty")},
		{"non-empty string", "test", nil},
		{"empty slice", []string{}, errors.New("value is empty")},
		{"non-empty slice", []string{"item"}, nil},
		{"zero integer", 0, errors.New("value is zero")},
		{"non-zero integer", 42, nil},
		{"false boolean", false, errors.New("value is false")},
		{"true boolean", true, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsNotEmpty(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

// TestIsAlphanumeric tests the IsAlphanumeric validator.
func TestIsAlphanumeric(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid alphanumeric", "abc123", nil},
		{"invalid characters", "abc@123", errors.New("value contains invalid characters")},
		{"not a string", 123, errors.New("value is not a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsAlphanumeric(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

// TestMinLength tests the MinLength validator.
func TestMinLength(t *testing.T) {
	minLength := MinLength(5)

	tests := []struct {
		name  string
		input string
		error error
	}{
		{"valid length", "testing", nil},
		{"too short", "test", errors.New("value must be at least 5 characters long")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := minLength(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

// TestMaxLength tests the MaxLength validator.
func TestMaxLength(t *testing.T) {
	maxLength := MaxLength(5)

	tests := []struct {
		name  string
		input string
		error error
	}{
		{"valid length", "test", nil},
		{"too long", "testing", errors.New("value must be at most 5 characters long")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := maxLength(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

// TestIsEmail tests the IsEmail validator.
func TestIsEmail(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid email", "user@example.com", nil},
		{"invalid email", "invalid-email", errors.New("value is not a valid email address")},
		{"not a string", 123, errors.New("value is not a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsEmail(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

// TestRegex tests the Regex validator.
func TestRegex(t *testing.T) {
	regex := Regex(`^\d+$`) // Matches strings with only digits

	tests := []struct {
		name  string
		input string
		error error
	}{
		{"valid regex", "12345", nil},
		{"invalid regex", "abc123", errors.New("value does not match the required pattern")},
		{"empty string", "", errors.New("value does not match the required pattern")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := regex(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

// TestValidate tests the Validate function.
func TestValidate(t *testing.T) {
	options := []ValidationOption{
		{
			Key:        "username",
			IsOptional: false,
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "Username is required"),
				CreateValidator(IsAlphanumeric, "Username must be alphanumeric"),
			},
		},
		{
			Key:        "email",
			IsOptional: false,
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "Email is required"),
				CreateValidator(IsEmail, "Invalid email address"),
			},
		},
		{
			Key:        "password",
			IsOptional: false,
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "Password is required"),
				CreateValidator(MinLength(6), "Password must be at least 6 characters"),
			},
		},
	}

	tests := []struct {
		name  string
		input map[string]interface{}
		error error
	}{
		{
			"valid input",
			map[string]interface{}{
				"username": "user123",
				"email":    "user@example.com",
				"password": "password123",
			},
			nil,
		},
		{
			"missing required field",
			map[string]interface{}{
				"username": "user123",
				"password": "password123",
			},
			errors.New("Email is required"),
		},
		{
			"invalid email",
			map[string]interface{}{
				"username": "user123",
				"email":    "invalid-email",
				"password": "password123",
			},
			errors.New("Invalid email address"),
		},
		{
			"password too short",
			map[string]interface{}{
				"username": "user123",
				"email":    "user@example.com",
				"password": "pass",
			},
			errors.New("Password must be at least 6 characters"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.input, options)
			require.Equal(t, test.error, err)
		})
	}
}
