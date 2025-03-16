package validator

import (
	"errors"
	"fmt"
	"reflect"
)

// MinLength checks if a string meets a minimum length requirement.
func MinLength(min int) ValidatorFunc {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value is not a string")
		}
		if len(str) < min {
			return fmt.Errorf("value must be at least %d characters long", min)
		}
		return nil
	}
}

// MaxLength checks if a string meets a maximum length requirement.
func MaxLength(max int) ValidatorFunc {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value is not a string")
		}
		if len(str) > max {
			return fmt.Errorf("value must be at most %d characters long", max)
		}
		return nil
	}
}

// Length checks if a string meets a length requirement within a range.
func Length(min, max int) ValidatorFunc {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value is not a string")
		}
		length := len(str)
		if length < min || length > max {
			return fmt.Errorf("value must be between %d and %d characters long", min, max)
		}
		return nil
	}
}

// MaxValue checks if a numeric value is less than or equal to a maximum value.
func Max(max float64) ValidatorFunc {
	return func(value interface{}) error {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if float64(v.Int()) > max {
				return fmt.Errorf("value must be less than or equal to %v", max)
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() > max {
				return fmt.Errorf("value must be less than or equal to %v", max)
			}
		default:
			return errors.New("value must be a number")
		}
		return nil
	}
}

// MinValue checks if a numeric value is greater than or equal to a minimum value.
func Min(min float64) ValidatorFunc {
	return func(value interface{}) error {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if float64(v.Int()) < min {
				return fmt.Errorf("value must be greater than or equal to %v", min)
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() < min {
				return fmt.Errorf("value must be greater than or equal to %v", min)
			}
		default:
			return errors.New("value must be a number")
		}
		return nil
	}
}

// Each checks if every element in a slice or array satisfies the provided validator function.
func Each(validatorFunc ValidatorFunc) ValidatorFunc {
	return func(value interface{}) error {
		// Check if the value is a slice or array
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return errors.New("value must be a slice or array")
		}

		// Iterate over each element and apply the validator function
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i).Interface()
			if err := validatorFunc(element); err != nil {
				return fmt.Errorf("element at index %d: %v", i, err)
			}
		}

		return nil
	}
}

// EachWithOptions applies a set of validation options to each element in a slice or array, returning the first error
func EachWithOptions(options []ValidationOption) ValidatorFunc {
	return func(value interface{}) error {
		if value == nil {
			return fmt.Errorf("value must be a non-nil slice or array")
		}
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return fmt.Errorf("value must be a slice or array, got %T", value)
		}
		if v.Len() == 0 {
			return nil
		}
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			nestedBody, ok := elem.(map[string]interface{})
			if !ok {
				if reflect.TypeOf(elem).Kind() == reflect.Struct {
					nestedBody = StructToMap(elem)
				} else {
					return fmt.Errorf("element at index %d must be an object, got %T", i, elem)
				}
			}
			if err := Validate(nestedBody, options); err != nil {
				return fmt.Errorf("element at index %d: %v", i, err)
			}
		}
		return nil
	}
}
