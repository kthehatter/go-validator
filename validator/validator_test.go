package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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
