package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestMaxValue(t *testing.T) {
	maxValue := Max(100)

	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid integer", 50, nil},
		{"valid float", 99.9, nil},
		{"equal to max", 100, nil},
		{"greater than max (integer)", 101, errors.New("value must be less than or equal to 100")},
		{"greater than max (float)", 100.1, errors.New("value must be less than or equal to 100")},
		{"invalid type (string)", "100", errors.New("value must be a number")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := maxValue(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestMinValue(t *testing.T) {
	minValue := Min(10)

	tests := []struct {
		name  string
		input interface{}
		error error
	}{
		{"valid integer", 50, nil},
		{"valid float", 10.1, nil},
		{"equal to min", 10, nil},
		{"less than min (integer)", 9, errors.New("value must be greater than or equal to 10")},
		{"less than min (float)", 9.9, errors.New("value must be greater than or equal to 10")},
		{"invalid type (string)", "10", errors.New("value must be a number")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := minValue(test.input)
			require.Equal(t, test.error, err)
		})
	}
}

func TestEach(t *testing.T) {
	// Test with IsString validator
	t.Run("Each element is a string", func(t *testing.T) {
		input := []interface{}{"hello", "world", "123"}
		err := Each(IsString)(input)
		require.NoError(t, err)
	})

	t.Run("Each element is not a string", func(t *testing.T) {
		input := []interface{}{"hello", 123, "world"}
		err := Each(IsString)(input)
		require.EqualError(t, err, "element at index 1: value must be a string")
	})

	// Test with IsNumber validator
	t.Run("Each element is a number", func(t *testing.T) {
		input := []interface{}{1, 2.5, 3}
		err := Each(IsNumber)(input)
		require.NoError(t, err)
	})

	t.Run("Each element is not a number", func(t *testing.T) {
		input := []interface{}{1, "2.5", 3}
		err := Each(IsNumber)(input)
		require.EqualError(t, err, "element at index 1: value must be a number")
	})

	// Test with IsArabic validator
	t.Run("Each element is Arabic", func(t *testing.T) {
		input := []interface{}{"مرحبا", "العالم", "١٢٣"}
		err := Each(IsArabic)(input)
		require.NoError(t, err)
	})

	t.Run("Each element is not Arabic", func(t *testing.T) {
		input := []interface{}{"مرحبا", "Hello", "العالم"}
		err := Each(IsArabic)(input)
		require.EqualError(t, err, "element at index 1: value must contain only Arabic characters")
	})

	// Test with non-slice/array input
	t.Run("Input is not a slice or array", func(t *testing.T) {
		input := "hello"
		err := Each(IsString)(input)
		require.EqualError(t, err, "value must be a slice or array")
	})
}

func TestEachWithOptions(t *testing.T) {
	nestedOptions := []ValidationOption{
		{
			Key:        "name",
			IsOptional: false,
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "name is required"),
				CreateValidator(IsString, "value must be a string"),
			},
		},
		{
			Key:        "value",
			IsOptional: false,
			Validators: []Validator{
				CreateValidator(IsNotEmpty, "value is required"),
			},
		},
	}

	tests := []struct {
		name     string
		input    interface{}
		expected error
	}{
		{
			name: "valid array of maps",
			input: []interface{}{
				map[string]interface{}{"name": "color", "value": "red"},
				map[string]interface{}{"name": "size", "value": "M"},
			},
			expected: nil,
		},
		{
			name: "invalid array - missing required field",
			input: []interface{}{
				map[string]interface{}{"name": "color", "value": "red"},
				map[string]interface{}{"value": "M"},
			},
			expected: errors.New("element at index 1: name is required"),
		},
		{
			name: "invalid array - wrong type",
			input: []interface{}{
				map[string]interface{}{"name": "color", "value": "red"},
				map[string]interface{}{"name": 123, "value": "M"},
			},
			expected: errors.New("element at index 1: value must be a string"),
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: nil,
		},
		{
			name:     "nil array",
			input:    nil,
			expected: errors.New("value must be a non-nil slice or array"),
		},
		{
			name:     "not an array",
			input:    "not an array",
			expected: errors.New("value must be a slice or array, got string"),
		},
		{
			name: "valid array of structs",
			input: []interface{}{
				struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				}{Name: "color", Value: "red"},
				struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				}{Name: "size", Value: "M"},
			},
			expected: nil,
		},
		{
			name: "invalid element type",
			input: []interface{}{
				"not a map or struct",
			},
			expected: errors.New("element at index 0 must be an object, got string"),
		},
		{
			name: "json with mixed case keys",
			input: []interface{}{
				map[string]interface{}{
					"name":  "ahmed",
					"Name":  123,
					"value": "test",
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := EachWithOptions(nestedOptions)
			err := validator(tt.input)
			if tt.expected == nil {
				require.NoError(t, err, "expected no error")
			} else {
				require.Error(t, err, "expected an error")
				require.Equal(t, tt.expected.Error(), err.Error(), "error message mismatch")
			}
		})
	}
}
