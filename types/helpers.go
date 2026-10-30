package types

import (
	"fmt"

	"github.com/victorrbravo/stories-service/constants"
)

// Nullable is a type that can have a special `null' value
type Nullable interface {
	isNil() bool
}

// MustBeNonEmptyDateTime checks that the user did not pass an empty date
func MustBeNonEmptyDateTime(date DateTime) error {
	if date.IsZero() {
		return fmt.Errorf(constants.DateCannotBeEmpty)
	}
	return nil
}

// AsInt convenience function to get a native `int' from our smart integer type
func AsInt(value Int) int {
	return value.get()
}

// AsString convenience function to get a native `string' from our smart string type
func AsString(value String) string {
	return value.get()
}

// AsBool convenience function to get a native `bool' from our smart boolean type
func AsBool(value Bool) bool {
	return value.get()
}

// IsNil evaluates whether a nullable thing is `nil'
func IsNil(value Nullable) bool {
	return value.isNil()
}
