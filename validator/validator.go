package validator

import "fmt"

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

// Validate checks the request body against the validation options and returns the first error.
func Validate(body map[string]interface{}, options []ValidationOption) error {
	for _, option := range options {
		value, exists := body[option.Key]

		// Skip validation if the field is optional and not present
		if option.IsOptional && !exists {
			continue
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
