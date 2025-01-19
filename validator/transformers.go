package validator

import (
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RemoveSpecialChars removes special characters from a string.
func RemoveSpecialChars(value interface{}) interface{} {
	if str, ok := value.(string); ok {
		var result strings.Builder
		for _, char := range str {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == ' ' {
				result.WriteRune(char)
			}
		}
		return result.String()
	}
	return value
}

// ToTitleCase converts a string to title case.
func ToTitleCase(value interface{}) interface{} {
	if str, ok := value.(string); ok {
		// Create a title caser for the English language
		caser := cases.Title(language.French)
		return caser.String(strings.ToLower(str))
	}
	return value
}

// ToInt converts a string or float to an integer.
func ToInt(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	case float64:
		return int(v)
	}
	return value
}

// ToFloat converts a string or integer to a float.
func ToFloat(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	case int:
		return float64(v)
	}
	return value
}

// Truncate truncates a string to a specified maximum length.
func Truncate(maxLength int) Transformer {
	return func(value interface{}) interface{} {
		if str, ok := value.(string); ok {
			if len(str) > maxLength {
				return str[:maxLength]
			}
		}
		return value
	}
}

// Replace replaces occurrences of a substring with another string.
func Replace(old, new string) Transformer {
	return func(value interface{}) interface{} {
		if str, ok := value.(string); ok {
			return strings.ReplaceAll(str, old, new)
		}
		return value
	}
}
