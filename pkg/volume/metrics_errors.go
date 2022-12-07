package volume

import (
	"fmt"
)

const (
	// ErrCodeNotSupported code for NotSupported Errors.
	ErrCodeNotSupported int = iota + 1
	ErrCodeNoPathDefined
	ErrCodeFsInfoFailed
)

// NewNotSupportedError creates a new MetricsError with code NotSupported.
func NewNotSupportedError() *MetricsError {
	return &MetricsError{
		Code: ErrCodeNotSupported,
		Msg:  "metrics are not supported for MetricsNil Volumes",
	}
}

// NewNotSupportedErrorWithDriverName creates a new MetricsError with code NotSupported.
// driver name is added to the error message.
func NewNotSupportedErrorWithDriverName(name string) *MetricsError {
	return &MetricsError{
		Code: ErrCodeNotSupported,
		Msg:  fmt.Sprintf("metrics are not supported for %s volumes", name),
	}
}

// NewNoPathDefinedError creates a new MetricsError with code NoPathDefined.
func NewNoPathDefinedError() *MetricsError {
	return &MetricsError{
		Code: ErrCodeNoPathDefined,
		Msg:  "no path defined for disk usage metrics.",
	}
}

// NewFsInfoFailedError creates a new MetricsError with code FsInfoFailed.
func NewFsInfoFailedError(err error) *MetricsError {
	return &MetricsError{
		Code: ErrCodeFsInfoFailed,
		Msg:  fmt.Sprintf("failed to get FsInfo due to error %v", err),
	}
}

// MetricsError to distinguish different Metrics Errors.
type MetricsError struct {
	Code int
	Msg  string
}

func (e *MetricsError) Error() string {
	return e.Msg
}

// IsNotSupported returns true if and only if err is "key" not found error.
func IsNotSupported(err error) bool {
	return isErrCode(err, ErrCodeNotSupported)
}

func isErrCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*MetricsError); ok {
		return e.Code == code
	}
	return false
}
