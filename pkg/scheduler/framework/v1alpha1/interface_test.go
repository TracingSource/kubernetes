package v1alpha1

import (
	"errors"
	"testing"
)

func TestStatus(t *testing.T) {
	tests := []struct {
		status            *Status
		expectedCode      Code
		expectedMessage   string
		expectedIsSuccess bool
		expectedAsError   error
	}{
		{
			status:            NewStatus(Success, ""),
			expectedCode:      Success,
			expectedMessage:   "",
			expectedIsSuccess: true,
			expectedAsError:   nil,
		},
		{
			status:            NewStatus(Error, "unknown error"),
			expectedCode:      Error,
			expectedMessage:   "unknown error",
			expectedIsSuccess: false,
			expectedAsError:   errors.New("unknown error"),
		},
		{
			status:            nil,
			expectedCode:      Success,
			expectedMessage:   "",
			expectedIsSuccess: true,
			expectedAsError:   nil,
		},
	}

	for i, test := range tests {
		if test.status.Code() != test.expectedCode {
			t.Errorf("test #%v, expect status.Code() returns %v, but %v", i, test.expectedCode, test.status.Code())
		}

		if test.status.Message() != test.expectedMessage {
			t.Errorf("test #%v, expect status.Message() returns %v, but %v", i, test.expectedMessage, test.status.Message())
		}

		if test.status.IsSuccess() != test.expectedIsSuccess {
			t.Errorf("test #%v, expect status.IsSuccess() returns %v, but %v", i, test.expectedIsSuccess, test.status.IsSuccess())
		}

		if test.status.AsError() == test.expectedAsError {
			continue
		}

		if test.status.AsError().Error() != test.expectedAsError.Error() {
			t.Errorf("test #%v, expect status.AsError() returns %v, but %v", i, test.expectedAsError, test.status.AsError())
		}
	}
}

// The String() method relies on the value and order of the status codes to function properly.
func TestStatusCodes(t *testing.T) {
	assertStatusCode(t, Success, 0)
	assertStatusCode(t, Error, 1)
	assertStatusCode(t, Unschedulable, 2)
	assertStatusCode(t, UnschedulableAndUnresolvable, 3)
	assertStatusCode(t, Wait, 4)
	assertStatusCode(t, Skip, 5)
}

func assertStatusCode(t *testing.T, code Code, value int) {
	if int(code) != value {
		t.Errorf("Status code %q should have a value of %v but got %v", code.String(), value, int(code))
	}
}
