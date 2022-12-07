package compiler

// Error represents compiler errors and their location in the document.
type Error struct {
	Context *Context
	Message string
}

// NewError creates an Error.
func NewError(context *Context, message string) *Error {
	return &Error{Context: context, Message: message}
}

// Error returns the string value of an Error.
func (err *Error) Error() string {
	if err.Context == nil {
		return "ERROR " + err.Message
	}
	return "ERROR " + err.Context.Description() + " " + err.Message
}

// ErrorGroup is a container for groups of Error values.
type ErrorGroup struct {
	Errors []error
}

// NewErrorGroupOrNil returns a new ErrorGroup for a slice of errors or nil if the slice is empty.
func NewErrorGroupOrNil(errors []error) error {
	if len(errors) == 0 {
		return nil
	} else if len(errors) == 1 {
		return errors[0]
	} else {
		return &ErrorGroup{Errors: errors}
	}
}

func (group *ErrorGroup) Error() string {
	result := ""
	for i, err := range group.Errors {
		if i > 0 {
			result += "\n"
		}
		result += err.Error()
	}
	return result
}
