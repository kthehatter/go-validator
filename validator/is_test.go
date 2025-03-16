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

func TestIsIn(t *testing.T) {
	isIn := IsIn("apple", "banana", "cherry")

	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid value", "apple", nil},
		{"invalid value", "grape", errors.New("value must be one of [apple banana cherry]")},
		{"wrong type", 123, errors.New("value must be one of [apple banana cherry]")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := isIn(test.input)
			require.Equal(t, test.error, err)
		})
	}
}
func TestIsNotIn(t *testing.T) {
	isNotIn := IsNotIn("apple", "banana", "cherry")
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid value", "grape", nil},
		{"invalid value", "apple", errors.New("value must not be one of [apple banana cherry]")},
		{"wrong type", 123, errors.New("value must not be one of [apple banana cherry]")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := isNotIn(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsInArray(t *testing.T) {
	isInArray := IsInArray([]string{"apple", "banana", "cherry"})
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid value", "apple", nil},
		{"invalid value", "grape", errors.New("value must be one of [apple banana cherry]")},
		{"wrong type", 123, errors.New("value must be one of [apple banana cherry]")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := isInArray(test.input)
			require.Equal(t, test.error, err)
		})
	}
}
func TestIsNotInArray(t *testing.T) {
	isNotInArray := IsNotInArray([]string{"apple", "banana", "cherry"})
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid value", "grape", nil},
		{"invalid value", "apple", errors.New("value must not be one of [apple banana cherry]")},
		{"wrong type", 123, errors.New("value must not be one of [apple banana cherry]")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := isNotInArray(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsString(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid string", "hello", nil},
		{"invalid type (int)", 123, errors.New("value must be a string")},
		{"invalid type (float)", 123.45, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsString(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid integer", 123, nil},
		{"valid float", 123.45, nil},
		{"invalid type (string)", "123", errors.New("value must be a number")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsNumber(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsInt(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid integer", 123, nil},
		{"invalid type (float)", 123.45, errors.New("value must be an integer")},
		{"invalid type (string)", "123", errors.New("value must be an integer")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsInt(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsFloat(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid float", 123.45, nil},
		{"invalid type (int)", 123, errors.New("value must be a float")},
		{"invalid type (string)", "123.45", errors.New("value must be a float")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsFloat(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsBool(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid bool (true)", true, nil},
		{"valid bool (false)", false, nil},
		{"invalid type (string)", "true", errors.New("value must be a boolean")},
		{"invalid type (int)", 1, errors.New("value must be a boolean")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsBool(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsSlice(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid slice", []int{1, 2, 3}, nil},
		{"invalid type (string)", "hello", errors.New("value must be a slice")},
		{"invalid type (int)", 123, errors.New("value must be a slice")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsSlice(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsMap(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid map", map[string]int{"a": 1, "b": 2}, nil},
		{"invalid type (string)", "hello", errors.New("value must be a map")},
		{"invalid type (int)", 123, errors.New("value must be a map")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsMap(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid URL", "https://example.com", nil},
		{"invalid URL", "example.com", errors.New("value is not a valid URL")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsURL(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid UUID", "123e4567-e89b-12d3-a456-426614174000", nil},
		{"invalid UUID", "invalid-uuid", errors.New("value is not a valid UUID")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsUUID(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsDate(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid date", "2023-10-01", nil},
		{"invalid date", "01-10-2023", errors.New("value is not a valid date (expected format: YYYY-MM-DD)")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsDate(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsTime(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid time", "15:04:05", nil},
		{"invalid time", "25:61:61", errors.New("value is not a valid time (expected format: HH:MM:SS)")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsTime(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsCreditCard(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid credit card", "4111 1111 1111 1111", nil},
		{"invalid credit card", "1234 5678 9012 3456", errors.New("value is not a valid credit card number")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsCreditCard(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsHexColor(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid hex color (6 digits)", "#FFFFFF", nil},
		{"valid hex color (3 digits)", "#FFF", nil},
		{"invalid hex color", "#ZZZ", errors.New("value is not a valid hexadecimal color code")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsHexColor(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid JSON", `{"key": "value"}`, nil},
		{"invalid JSON", `{"key": "value"`, errors.New("value is not valid JSON")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsJSON(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsIP(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid IPv4", "192.168.1.1", nil},
		{"valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", nil},
		{"invalid IP", "invalid-ip", errors.New("value is not a valid IP address")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsIP(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid alpha", "Hello", nil},
		{"invalid alpha", "Hello123", errors.New("value must contain only alphabetic characters")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsAlpha(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid alphanumeric", "Hello123", nil},
		{"invalid alphanumeric", "Hello@123", errors.New("value must contain only alphanumeric characters")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsAlphaNumeric(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsBase64(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid Base64", "aGVsbG8=", nil}, // "hello" in Base64
		{"invalid Base64", "aGVsbG8", errors.New("value is not valid Base64")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsBase64(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsArabic(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid Arabic", "مرحبا", nil},
		{"invalid Arabic", "Hello123", errors.New("value must contain only Arabic characters")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsArabic(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestIsAlphaArabic(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid Arabic and Latin", "مرحبا Hello", nil},
		{"invalid Arabic and Latin", "مرحبا123", errors.New("value must contain only Arabic and Latin alphabetic characters")},
		{"invalid type (int)", 123, errors.New("value must be a string")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsAlphaArabic(test.input)
			require.Equal(t, test.error, err)
		})
	}
}
func TestIsBase64Image(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{
			name:  "valid Base64 image",
			input: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
			error: nil,
		},
		{
			name:  "invalid Base64 image",
			input: "data:image/png;base64,invalid",
			error: errors.New("value is not valid Base64 image"),
		},
		{
			name:  "invalid type (int)",
			input: 123,
			error: errors.New("value must be a string"),
		},
		{
			name:  "invalid base64 image format (missing prefix)",
			input: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wcAAgAB/1h1ZAAAAABJRU5ErkJggg==",
			error: errors.New("invalid base64 image format: must start with 'data:image/'"),
		},
		{
			name:  "unsupported image format",
			input: "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjwvc3ZnPg==",
			error: errors.New("invalid image format"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := IsBase64Image(test.input)
			require.Equal(t, test.error, err)
		})
	}
}
