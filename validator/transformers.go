package validator

import (
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// applyToArrayOrValue is a helper function to handle both single values and arrays
func applyToArrayOrValue(value any, transform func(any) any) any {
	// Use reflection to check if the value is a slice/array
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		// Create a new slice to store transformed values
		result := make([]any, v.Len())
		for i := range v.Len() {
			result[i] = transform(v.Index(i).Interface())
		}
		return result
	}
	// If not a slice, apply transformation directly
	return transform(value)
}

// ToLower transforms a string or array of strings to lowercase
func ToLower(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		if str, ok := v.(string); ok {
			return strings.ToLower(str)
		}
		return v
	})
}

// ToUpper transforms a string or array of strings to uppercase
func ToUpper(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		if str, ok := v.(string); ok {
			return strings.ToUpper(str)
		}
		return v
	})
}

// Trim trims leading and trailing whitespace from a string or array of strings
func Trim(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		if str, ok := v.(string); ok {
			return strings.TrimSpace(str)
		}
		return v
	})
}

// RemoveSpecialChars removes special characters from a string or array of strings
func RemoveSpecialChars(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		if str, ok := v.(string); ok {
			var result strings.Builder
			for _, char := range str {
				if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == ' ' {
					result.WriteRune(char)
				}
			}
			return result.String()
		}
		return v
	})
}

// ToTitleCase converts a string or array of strings to title case
func ToTitleCase(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		if str, ok := v.(string); ok {
			caser := cases.Title(language.French)
			return caser.String(strings.ToLower(str))
		}
		return v
	})
}

// ToInt converts a string/float or array of strings/floats to integer(s)
func ToInt(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		switch val := v.(type) {
		case string:
			if i, err := strconv.Atoi(val); err == nil {
				return i
			}
		case float64:
			return int(val)
		}
		return v
	})
}

// ToFloat converts a string/integer or array of strings/integers to float(s)
func ToFloat(value any) any {
	return applyToArrayOrValue(value, func(v any) any {
		switch val := v.(type) {
		case string:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				return f
			}
		case int:
			return float64(val)
		}
		return v
	})
}

// Truncate truncates a string or array of strings to a specified maximum length
func Truncate(maxLength int) Transformer {
	return func(value any) any {
		return applyToArrayOrValue(value, func(v any) any {
			if str, ok := v.(string); ok {
				if len(str) > maxLength {
					return str[:maxLength]
				}
			}
			return v
		})
	}
}

// Replace replaces occurrences of a substring with another string in a string or array of strings
func Replace(old, new string) Transformer {
	return func(value any) any {
		return applyToArrayOrValue(value, func(v any) any {
			if str, ok := v.(string); ok {
				return strings.ReplaceAll(str, old, new)
			}
			return v
		})
	}
}
