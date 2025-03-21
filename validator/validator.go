package validator

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidatorFunc is a function that validates a field and returns an error if validation fails.
type ValidatorFunc func(value interface{}) error

// TransformerFunc is a function that transforms a value.
type Transformer func(value interface{}) interface{}

// Validator defines a validator function and its error message.
type Validator struct {
	Func    ValidatorFunc
	Message string
}

// ValidationOption defines the validation rules for a specific field.
type ValidationOption struct {
	Key          string             // Field name in the request body
	IsOptional   bool               // Whether the field is optional
	Validators   []Validator        // List of validators for the field
	Transformers []Transformer      // List of transformers for the field
	Nested       []ValidationOption // Validation options for nested objects
}

// Validate validates the request body based on the provided validation options.
// It applies any transformers to update field values, enforces required fields,
// and executes associated validators, returning the first error encountered.
// For fields with nested validation options, it recursively validates the sub-objects.
func Validate(body map[string]interface{}, options []ValidationOption) error {
	for _, option := range options {
		value, exists := body[option.Key]

		// Skip validation if the field is optional and not present
		if option.IsOptional && !exists {
			continue
		}

		// Check if the field is required but missing *before* running validators
		if !option.IsOptional && !exists {
			return fmt.Errorf("%s is required", option.Key)
		}

		// Apply transformations
		if exists {
			for _, transformer := range option.Transformers {
				value = transformer(value)
			}
			body[option.Key] = value // Update the body with the transformed value
		}

		// Run all validators for the field (if it exists)
		for _, validator := range option.Validators {
			if err := validator.Func(value); err != nil {
				if validator.Message == "" {
					return err
				}
				return fmt.Errorf("%s", validator.Message)
			}
		}

		// Check if the field is required but missing
		if !option.IsOptional && !exists {
			return fmt.Errorf("'%s' is required", option.Key)
		}

		// Handle nested validation
		if option.Nested != nil {
			nestedBody, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("'%s' must be an object", option.Key)
			}
			if err := Validate(nestedBody, option.Nested); err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to create a validator
func CreateValidator(fn ValidatorFunc, message string) Validator {
	return Validator{
		Func:    fn,
		Message: message,
	}
}

// StructToMap converts a struct to a map[string]interface{} for validation
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := range v.NumField() {
		field := t.Field(i)
		// Use the json tag if present, otherwise fall back to field name
		jsonTag := field.Tag.Get("json")
		key := field.Name
		if jsonTag != "" {
			// Split on "," to handle options like "omitempty"
			if parts := strings.Split(jsonTag, ","); len(parts) > 0 {
				key = parts[0]
			}
		}
		// Only process exported fields
		if field.PkgPath == "" { // PkgPath is empty for exported fields
			result[key] = v.Field(i).Interface()
		}
	}
	return result
}
