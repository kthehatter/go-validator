package validator

import (
	"reflect"
	"testing"
)

func TestToLower(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"Hello World", "hello world"}, // Single string
		{"HELLO WORLD", "hello world"}, // Single string
		{"123ABC", "123abc"},           // Mixed alphanumeric
		{"", ""},                       // Empty string
		{123, 123},                     // Non-string input
		{nil, nil},                     // Nil input
		{[]any{"Hello", "WORLD"}, []any{"hello", "world"}}, // Array of strings
		{[]string{"ABC", "XYZ"}, []any{"abc", "xyz"}},      // Slice of strings
		{[]any{"Hello", 123}, []any{"hello", 123}},         // Mixed array
	}

	for _, test := range tests {
		result := ToLower(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ToLower(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"Hello World", "HELLO WORLD"}, // Single string
		{"hello world", "HELLO WORLD"}, // Single string
		{"123abc", "123ABC"},           // Mixed alphanumeric
		{"", ""},                       // Empty string
		{123, 123},                     // Non-string input
		{nil, nil},                     // Nil input
		{[]any{"Hello", "world"}, []any{"HELLO", "WORLD"}}, // Array of strings
		{[]string{"abc", "xyz"}, []any{"ABC", "XYZ"}},      // Slice of strings
		{[]any{"hello", 456}, []any{"HELLO", 456}},         // Mixed array
	}

	for _, test := range tests {
		result := ToUpper(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ToUpper(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"  Hello World  ", "Hello World"}, // Single string with spaces
		{"\tHello World\n", "Hello World"}, // Single string with tabs/newlines
		{"Hello World", "Hello World"},     // Single string no trimming needed
		{"", ""},                           // Empty string
		{123, 123},                         // Non-string input
		{nil, nil},                         // Nil input
		{[]any{"  Hello  ", "\tWorld\n"}, []any{"Hello", "World"}}, // Array of strings
		{[]string{"  abc  ", " xyz "}, []any{"abc", "xyz"}},        // Slice of strings
	}

	for _, test := range tests {
		result := Trim(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Trim(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestRemoveSpecialChars(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"Hello!@# World!", "Hello World"}, // Single string with special chars
		{"123#456", "123456"},              // Single string with numbers
		{"", ""},                           // Empty string
		{123, 123},                         // Non-string input
		{[]any{"Hello!@#", "World$%"}, []any{"Hello", "World"}}, // Array of strings
		{[]string{"abc#@!", "123$%"}, []any{"abc", "123"}},      // Slice of strings
	}

	for _, test := range tests {
		result := RemoveSpecialChars(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("RemoveSpecialChars(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"hello world", "Hello World"},     // Single string
		{"HELLO WORLD", "Hello World"},     // Single string uppercase
		{"çimant", "Çimant"},               // Single string with special chars
		{"hello, world!", "Hello, World!"}, // Single string with punctuation
		{"123", "123"},                     // Single string numeric
		{123, 123},                         // Non-string input
		{[]any{"hello world", "GOOD BYE"}, []any{"Hello World", "Good Bye"}}, // Array of strings
		{[]string{"abc def", "xyz"}, []any{"Abc Def", "Xyz"}},                // Slice of strings
	}

	for _, test := range tests {
		result := ToTitleCase(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ToTitleCase(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"123", 123},                             // Single string
		{123.45, 123},                            // Single float
		{"abc", "abc"},                           // Invalid single string
		{[]any{"123", "456"}, []any{123, 456}},   // Array of valid strings
		{[]string{"789", "0"}, []any{789, 0}},    // Slice of valid strings
		{[]any{"abc", "123"}, []any{"abc", 123}}, // Mixed array
	}

	for _, test := range tests {
		result := ToInt(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ToInt(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{"123.45", 123.45}, // Single string
		{123, 123.0},       // Single int
		{"abc", "abc"},     // Invalid single string
		{[]any{"123.45", "0"}, []any{123.45, 0.0}},    // Array of valid strings
		{[]string{"789", "1.23"}, []any{789.0, 1.23}}, // Slice of valid strings
		{[]any{"abc", "123"}, []any{"abc", 123.0}},    // Mixed array
	}

	for _, test := range tests {
		result := ToFloat(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ToFloat(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	transformer := Truncate(5)
	tests := []struct {
		input    any
		expected any
	}{
		{"hello world", "hello"}, // Single string
		{"hi", "hi"},             // Single string shorter than max
		{123, 123},               // Non-string input
		{[]any{"hello world", "abcdef"}, []any{"hello", "abcde"}}, // Array of strings
		{[]string{"hi", "goodbye"}, []any{"hi", "goodb"}},         // Slice of strings
	}

	for _, test := range tests {
		result := transformer(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Truncate(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestReplace(t *testing.T) {
	transformer := Replace("foo", "bar")
	tests := []struct {
		input    any
		expected any
	}{
		{"foo bar", "bar bar"},                                    // Single string
		{"foo bar foo", "bar bar bar"},                            // Single string multiple replacements
		{"hello world", "hello world"},                            // Single string no replacement
		{123, 123},                                                // Non-string input
		{[]any{"foo", "foo bar"}, []any{"bar", "bar bar"}},        // Array of strings
		{[]string{"hello foo", "foo"}, []any{"hello bar", "bar"}}, // Slice of strings
	}

	for _, test := range tests {
		result := transformer(test.input)
		if !reflect.DeepEqual(result, test.expected) {
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
		input       map[string]any
		expectedErr string
	}{
		{
			name: "Valid input",
			input: map[string]any{
				"username": "  JohnDoe  ",
				"email":    "  JOHN@EXAMPLE.COM  ",
			},
			expectedErr: "", // No error expected
		},
		{
			name: "Username too short",
			input: map[string]any{
				"username": "  Jo  ", // After trimming and lowercasing: "jo" (length 2)
				"email":    "john@example.com",
			},
			expectedErr: "username must be at least 3 characters long",
		},
		{
			name: "Empty username",
			input: map[string]any{
				"username": "  ", // After trimming: "" (empty)
				"email":    "john@example.com",
			},
			expectedErr: "username is required",
		},
		{
			name: "Invalid email",
			input: map[string]any{
				"username": "johndoe",
				"email":    "invalid-email", // Invalid email format
			},
			expectedErr: "invalid email address",
		},
		{
			name: "Missing required field",
			input: map[string]any{
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
