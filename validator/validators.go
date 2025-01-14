package validator

import (
	"errors"
	"fmt"
	"regexp"
)

// Regex validates a string against a regular expression.
func Regex(pattern string) ValidatorFunc {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("Invalid regex pattern: %s", err))
	}
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("value must be a string")
		}

		if !re.MatchString(str) {
			return errors.New("value does not match the required pattern")
		}

		return nil
	}
}
