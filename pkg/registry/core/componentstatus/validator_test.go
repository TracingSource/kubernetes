package componentstatus

import (
	"errors"
	"fmt"
	"testing"

	"k8s.io/kubernetes/pkg/probe"
)

func matchError(data []byte) error {
	if string(data) != "bar" {
		return errors.New("match error")
	}
	return nil
}

func TestValidate(t *testing.T) {
	tests := []struct {
		probeResult probe.Result
		probeData   string
		probeErr    error

		expectResult probe.Result
		expectData   string
		expectErr    bool

		validator ValidatorFn
	}{
		{probe.Unknown, "", fmt.Errorf("probe error"), probe.Unknown, "", true, nil},
		{probe.Failure, "", nil, probe.Failure, "", false, nil},
		{probe.Success, "foo", nil, probe.Failure, "foo", true, matchError},
		{probe.Success, "foo", nil, probe.Success, "foo", false, nil},
	}

	s := Server{Addr: "foo.com", Port: 8080, Path: "/healthz"}

	for _, test := range tests {
		fakeProber := &fakeHttpProber{
			result: test.probeResult,
			body:   test.probeData,
			err:    test.probeErr,
		}

		s.Validate = test.validator
		s.Prober = fakeProber
		result, data, err := s.DoServerCheck()
		if test.expectErr && err == nil {
			t.Error("unexpected non-error")
		}
		if !test.expectErr && err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if data != test.expectData {
			t.Errorf("expected %s, got %s", test.expectData, data)
		}
		if result != test.expectResult {
			t.Errorf("expected %s, got %s", test.expectResult, result)
		}
	}
}
