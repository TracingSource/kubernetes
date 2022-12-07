package validation

import (
	"fmt"
)

// Error is the type that's returned when the validation of an APIs arguments constraints fails.
type Error struct {
	// PackageType is the package type of the object emitting the error. For types, the value
	// matches that produced the the '%T' format specifier of the fmt package. For other elements,
	// such as functions, it is just the package name (e.g., "autorest").
	PackageType string

	// Method is the name of the method raising the error.
	Method string

	// Message is the error message.
	Message string
}

// Error returns a string containing the details of the validation failure.
func (e Error) Error() string {
	return fmt.Sprintf("%s#%s: Invalid input: %s", e.PackageType, e.Method, e.Message)
}

// NewError creates a new Error object with the specified parameters.
// message is treated as a format string to which the optional args apply.
func NewError(packageType string, method string, message string, args ...interface{}) Error {
	return Error{
		PackageType: packageType,
		Method:      method,
		Message:     fmt.Sprintf(message, args...),
	}
}
