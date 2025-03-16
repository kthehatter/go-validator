package validator

import (
	"testing"
)

func TestToLower(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"Hello World", "hello world"}, // Lowercase conversion
		{"HELLO WORLD", "hello world"}, // Lowercase conversion
		{"123ABC", "123abc"},           // Mixed alphanumeric
		{"", ""},                       // Empty string
		{123, 123},                     // Non-string input
		{nil, nil},                     // Nil input
	}

	for _, test := range tests {
		result := ToLower(test.input)
		if result != test.expected {
			t.Errorf("ToLower(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"Hello World", "HELLO WORLD"}, // Uppercase conversion
		{"hello world", "HELLO WORLD"}, // Uppercase conversion
		{"123abc", "123ABC"},           // Mixed alphanumeric
		{"", ""},                       // Empty string
		{123, 123},                     // Non-string input
		{nil, nil},                     // Nil input
	}

	for _, test := range tests {
		result := ToUpper(test.input)
		if result != test.expected {
			t.Errorf("ToUpper(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"  Hello World  ", "Hello World"}, // Trim leading and trailing spaces
		{"\tHello World\n", "Hello World"}, // Trim tabs and newlines
		{"Hello World", "Hello World"},     // No trimming needed
		{"", ""},                           // Empty string
		{123, 123},                         // Non-string input
		{nil, nil},                         // Nil input
	}

	for _, test := range tests {
		result := Trim(test.input)
		if result != test.expected {
			t.Errorf("Trim(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestRemoveSpecialChars(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"Hello!@# World!", "Hello World"},
		{"123#456", "123456"},
		{123, 123}, // Non-string input
	}

	for _, test := range tests {
		result := RemoveSpecialChars(test.input)
		if result != test.expected {
			t.Errorf("RemoveSpecialChars(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"hello world", "Hello World"},
		{"HELLO WORLD", "Hello World"},
		{"çimant", "Çimant"},
		{"hello, world!", "Hello, World!"}, // Handles punctuation correctly
		{"123", "123"},                     // Non-alphabetic input
		{123, 123},                         // Non-string input
		{"السلام عليكم", "السلام عليكم"}, // Arabic input
	}

	for _, test := range tests {
		result := ToTitleCase(test.input)
		if result != test.expected {
			t.Errorf("ToTitleCase(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"123", 123},
		{123.45, 123},
		{"abc", "abc"}, // Invalid input
	}

	for _, test := range tests {
		result := ToInt(test.input)
		if result != test.expected {
			t.Errorf("ToInt(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"123.45", 123.45},
		{123, 123.0},
		{"abc", "abc"}, // Invalid input
	}

	for _, test := range tests {
		result := ToFloat(test.input)
		if result != test.expected {
			t.Errorf("ToFloat(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	transformer := Truncate(5)

	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"hello world", "hello"},
		{"hi", "hi"},
		{123, 123}, // Non-string input
	}

	for _, test := range tests {
		result := transformer(test.input)
		if result != test.expected {
			t.Errorf("Truncate(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestReplace(t *testing.T) {
	transformer := Replace("foo", "bar")

	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"foo bar", "bar bar"},
		{"foo bar foo", "bar bar bar"},
		{"hello world", "hello world"},
		{123, 123}, // Non-string input
	}

	for _, test := range tests {
		result := transformer(test.input)
		if result != test.expected {
			t.Errorf("Replace(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestValidationWithTransformations(t *testing.T) {
	// Define validation rules with transformations
	rules := []ValidationOption{
		{
			Key:        "username",
			IsOptional: false,
			Transformers: []Transformer{
				Trim,
				ToLower,
			},
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "username is required"),
				CreateValidator(MinLength(3), "username must be at least 3 characters long"),
			},
		},
		{
			Key:        "email",
			IsOptional: false,
			Transformers: []Transformer{
				Trim,
				ToLower,
			},
			Validators: []Validator{
				CreateValidator(IsEmail, "invalid email address"),
			},
		},
	}

	// Test cases
	tests := []struct {
		name        string
		input       map[string]interface{}
		expectedErr string
	}{
		{
			name: "Valid input",
			input: map[string]interface{}{
				"username": "  JohnDoe  ",
				"email":    "  JOHN@EXAMPLE.COM  ",
			},
			expectedErr: "", // No error expected
		},
		{
			name: "Username too short",
			input: map[string]interface{}{
				"username": "  Jo  ", // After trimming and lowercasing: "jo" (length 2)
				"email":    "john@example.com",
			},
			expectedErr: "username must be at least 3 characters long",
		},
		{
			name: "Empty username",
			input: map[string]interface{}{
				"username": "  ", // After trimming: "" (empty)
				"email":    "john@example.com",
			},
			expectedErr: "username is required",
		},
		{
			name: "Invalid email",
			input: map[string]interface{}{
				"username": "johndoe",
				"email":    "invalid-email", // Invalid email format
			},
			expectedErr: "invalid email address",
		},
		{
			name: "Missing required field",
			input: map[string]interface{}{
				"email": "john@example.com", // Missing "username"
			},
			expectedErr: "username is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.input, rules)
			if test.expectedErr == "" {
				// No error expected
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				// Error expected
				if err == nil {
					t.Errorf("Expected error: %v, got nil", test.expectedErr)
				} else if err.Error() != test.expectedErr {
					t.Errorf("Expected error: %v, got: %v", test.expectedErr, err.Error())
				}
			}
		})
	}
}
