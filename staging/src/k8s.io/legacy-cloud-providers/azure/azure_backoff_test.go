// +build !providerless

package azure

import (
	"fmt"
	"net/http"
	"testing"
)

func TestShouldRetryHTTPRequest(t *testing.T) {
	tests := []struct {
		code     int
		err      error
		expected bool
	}{
		{
			code:     http.StatusBadRequest,
			expected: true,
		},
		{
			code:     http.StatusInternalServerError,
			expected: true,
		},
		{
			code:     http.StatusOK,
			err:      fmt.Errorf("some error"),
			expected: true,
		},
		{
			code:     http.StatusOK,
			expected: false,
		},
		{
			code:     399,
			expected: false,
		},
	}
	for _, test := range tests {
		resp := &http.Response{
			StatusCode: test.code,
		}
		res := shouldRetryHTTPRequest(resp, test.err)
		if res != test.expected {
			t.Errorf("expected: %v, saw: %v", test.expected, res)
		}
	}
}

func TestIsSuccessResponse(t *testing.T) {
	tests := []struct {
		code     int
		expected bool
	}{
		{
			code:     http.StatusNotFound,
			expected: false,
		},
		{
			code:     http.StatusInternalServerError,
			expected: false,
		},
		{
			code:     http.StatusOK,
			expected: true,
		},
	}

	for _, test := range tests {
		resp := http.Response{
			StatusCode: test.code,
		}
		res := isSuccessHTTPResponse(&resp)
		if res != test.expected {
			t.Errorf("expected: %v, saw: %v", test.expected, res)
		}
	}
}

func TestProcessRetryResponse(t *testing.T) {
	az := &Cloud{}
	tests := []struct {
		code int
		err  error
		stop bool
	}{
		{
			code: http.StatusBadRequest,
			stop: false,
		},
		{
			code: http.StatusInternalServerError,
			stop: false,
		},
		{
			code: http.StatusSeeOther,
			err:  fmt.Errorf("some error"),
			stop: false,
		},
		{
			code: http.StatusSeeOther,
			stop: true,
		},
		{
			code: http.StatusOK,
			stop: true,
		},
		{
			code: http.StatusOK,
			err:  fmt.Errorf("some error"),
			stop: false,
		},
		{
			code: 399,
			stop: true,
		},
	}

	for _, test := range tests {
		resp := &http.Response{
			StatusCode: test.code,
		}
		res, err := az.processHTTPRetryResponse(nil, "", resp, test.err)
		if res != test.stop {
			t.Errorf("expected: %v, saw: %v", test.stop, res)
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
