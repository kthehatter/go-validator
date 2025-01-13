package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsNotEmpty tests the IsNotEmpty validator.
func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected error
	}{
		// Strings
		{name: "Non-empty string", value: "hello", expected: nil},
		{name: "Empty string", value: "", expected: errors.New("value is empty")},

		// Numbers
		{name: "Non-zero int", value: 42, expected: nil},
		{name: "Zero int", value: 0, expected: errors.New("value is zero")},
		{name: "Non-zero float", value: 3.14, expected: nil},
		{name: "Zero float", value: 0.0, expected: errors.New("value is zero")},

		// Booleans
		{name: "True boolean", value: true, expected: nil},
		{name: "False boolean", value: false, expected: errors.New("value is false")},

		// Slices
		{name: "Non-empty slice", value: []int{1, 2, 3}, expected: nil},
		{name: "Empty slice", value: []int{}, expected: errors.New("value is empty")},

		// Maps
		{name: "Non-empty map", value: map[string]int{"a": 1}, expected: nil},
		{name: "Empty map", value: map[string]int{}, expected: errors.New("value is empty")},

		// Nil
		{name: "Nil value", value: nil, expected: errors.New("value is nil")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsNotEmpty(tt.value)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
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

func TestNestedValidation(t *testing.T) {
	options := []ValidationOption{
		{
			Key: "user",
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "User is required"),
			},
			Nested: []ValidationOption{
				{
					Key: "name",
					Validators: []Validator{
						CreateValidator(IsNotEmpty, "Name is required"),
					},
				},
				{
					Key: "age",
					Validators: []Validator{
						CreateValidator(IsNotEmpty, "age is required"),
						CreateValidator(func(value interface{}) error {
							age, ok := value.(float64)
							if !ok {
								return errors.New("age must be a number")
							}
							if age < 0 {
								return errors.New("age must be positive")
							}
							return errors.New("Invalid age")
						}, "Invalid age"),
					},
				},
			},
		},
	}

	tests := []struct {
		name  string
		input map[string]interface{}
		error error
	}{
		{
			"valid nested input",
			map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John Doe",
					"age":  30,
				},
			},
			errors.New("Invalid age"),
		},
		{
			"missing nested field",
			map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John Doe",
				},
			},
			errors.New("age is required"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.input, options)
			require.Equal(t, test.error, err)
		})
	}
}

func TestCustomValidator(t *testing.T) {
	customValidator := CreateValidator(func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value must be a string")
		}
		if str != "expected" {
			return errors.New("Value must be 'expected'")
		}
		return nil
	}, "Value must be 'expected'")

	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid input", "expected", nil},
		{"invalid input", "unexpected", errors.New("Value must be 'expected'")},
		{"wrong type", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := customValidator.Func(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIntegration(t *testing.T) {
	options := []ValidationOption{
		{
			Key: "username",
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "Username is required"),
				CreateValidator(IsAlphanumeric, "Username must be alphanumeric"),
			},
		},
		{
			Key: "email",
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "Email is required"),
				CreateValidator(IsEmail, "Invalid email address"),
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
			},
			nil,
		},
		{
			"invalid username",
			map[string]interface{}{
				"username": "user@123",
				"email":    "user@example.com",
			},
			errors.New("Username must be alphanumeric"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.input, options)
			require.Equal(t, test.error, err)
		})
	}
}
