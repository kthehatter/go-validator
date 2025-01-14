package validator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IsNotEmpty checks if a value is not empty.
func IsNotEmpty(value interface{}) error {
	if value == nil {
		return errors.New("value is nil")
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return errors.New("value is empty")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return errors.New("value is zero")
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return errors.New("value is zero")
		}
	case reflect.Bool:
		if !v.Bool() {
			return errors.New("value is false")
		}
	default:
		// For unsupported types, assume the value is not empty
		return nil
	}

	return nil
}

// IsAlphanumeric checks if a string contains only alphanumeric characters.
func IsAlphanumeric(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}
	for _, char := range str {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return errors.New("value contains invalid characters")
		}
	}
	return nil
}

// IsEmail checks if a string is a valid email address.
func IsEmail(value interface{}) error {
	// Check if the input is a string
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	// Use net/mail to validate the email
	_, err := mail.ParseAddress(str)
	if err != nil {
		return errors.New("value is not a valid email address")
	}

	return nil
}

// IsIn checks if a value is in a predefined list of allowed values.
func IsIn(allowedValues ...interface{}) ValidatorFunc {
	return func(value interface{}) error {
		for _, allowed := range allowedValues {
			if value == allowed {
				return nil
			}
		}
		return fmt.Errorf("value must be one of %v", allowedValues)
	}
}

// IsString checks if a value is a string.
func IsString(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.String {
		return errors.New("value must be a string")
	}
	return nil
}

// IsNumber checks if a value is a number (int or float).
func IsNumber(value interface{}) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return nil
	default:
		return errors.New("value must be a number")
	}
}

// IsInt checks if a value is an integer.
func IsInt(value interface{}) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return nil
	default:
		return errors.New("value must be an integer")
	}
}

// IsFloat checks if a value is a float.
func IsFloat(value interface{}) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return nil
	default:
		return errors.New("value must be a float")
	}
}

// IsBool checks if a value is a boolean.
func IsBool(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Bool {
		return errors.New("value must be a boolean")
	}
	return nil
}

// IsSlice checks if a value is a slice.
func IsSlice(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		return errors.New("value must be a slice")
	}
	return nil
}

// IsMap checks if a value is a map.
func IsMap(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return errors.New("value must be a map")
	}
	return nil
}

// IsURL checks if a string is a valid URL.
func IsURL(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return errors.New("value is not a valid URL")
	}
	return nil
}

// IsUUID checks if a string is a valid UUID.
func IsUUID(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	uuidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	matched, err := regexp.MatchString(uuidRegex, str)
	if err != nil || !matched {
		return errors.New("value is not a valid UUID")
	}
	return nil
}

// IsDate checks if a string is a valid date in the format YYYY-MM-DD.
func IsDate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		return errors.New("value is not a valid date (expected format: YYYY-MM-DD)")
	}
	return nil
}

// IsTime checks if a string is a valid time in the format HH:MM:SS.
func IsTime(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	_, err := time.Parse("15:04:05", str)
	if err != nil {
		return errors.New("value is not a valid time (expected format: HH:MM:SS)")
	}
	return nil
}

// IsCreditCard checks if a string is a valid credit card number using the Luhn algorithm.
func IsCreditCard(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	// Remove spaces and dashes
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "-", "")

	// Check if the string is a valid number
	if _, err := strconv.Atoi(str); err != nil {
		return errors.New("value is not a valid credit card number")
	}

	// Luhn algorithm
	sum := 0
	alternate := false
	for i := len(str) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(str[i]))
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum += digit
		alternate = !alternate
	}

	if sum%10 != 0 {
		return errors.New("value is not a valid credit card number")
	}
	return nil
}

// IsHexColor checks if a string is a valid hexadecimal color code.
func IsHexColor(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	hexColorRegex := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	matched, err := regexp.MatchString(hexColorRegex, str)
	if err != nil || !matched {
		return errors.New("value is not a valid hexadecimal color code")
	}
	return nil
}

// IsJSON checks if a string is valid JSON.
func IsJSON(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	var js json.RawMessage
	if err := json.Unmarshal([]byte(str), &js); err != nil {
		return errors.New("value is not valid JSON")
	}
	return nil
}

// IsIP checks if a string is a valid IP address (IPv4 or IPv6).
func IsIP(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	if net.ParseIP(str) == nil {
		return errors.New("value is not a valid IP address")
	}
	return nil
}

// IsAlpha checks if a string contains only alphabetic characters.
func IsAlpha(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	alphaRegex := `^[a-zA-Z]+$`
	matched, err := regexp.MatchString(alphaRegex, str)
	if err != nil || !matched {
		return errors.New("value must contain only alphabetic characters")
	}
	return nil
}

// IsAlphaNumeric checks if a string contains only alphanumeric characters.
func IsAlphaNumeric(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	alphaNumericRegex := `^[a-zA-Z0-9]+$`
	matched, err := regexp.MatchString(alphaNumericRegex, str)
	if err != nil || !matched {
		return errors.New("value must contain only alphanumeric characters")
	}
	return nil
}

// IsArabic checks if a string contains only Arabic characters (including spaces and common Arabic punctuation).
func IsArabic(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}

	// Regex to match Arabic characters (Unicode block for Arabic and Arabic Supplement)
	arabicRegex := `^[\p{Arabic}\s]+$`
	matched, err := regexp.MatchString(arabicRegex, str)
	if err != nil {
		return errors.New("an error occurred while validating the string")
	}
	if !matched {
		return errors.New("value must contain only Arabic characters")
	}
	return nil
}

// IsAlphaArabic checks if a string contains only Arabic and Latin alphabetic characters.
func IsAlphaArabic(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	// Regex to match Arabic and Latin alphabetic characters
	alphaArabicRegex := `^[\p{Arabic}\p{Latin}\s]+$`
	matched, err := regexp.MatchString(alphaArabicRegex, str)
	if err != nil || !matched {
		return errors.New("value must contain only Arabic and Latin alphabetic characters")
	}
	return nil
}

// IsBase64 checks if a string is valid Base64-encoded data.
func IsBase64(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be a string")
	}
	_, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.New("value is not valid Base64")
	}
	return nil
}
