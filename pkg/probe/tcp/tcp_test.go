package tcp

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"k8s.io/kubernetes/pkg/probe"
)

func TestTcpHealthChecker(t *testing.T) {
	// Setup a test server that responds to probing correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	tHost, tPortStr, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	tPort, err := strconv.Atoi(tPortStr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tests := []struct {
		host string
		port int

		expectedStatus probe.Result
		expectedError  error
	}{
		// A connection is made and probing would succeed
		{tHost, tPort, probe.Success, nil},
		// No connection can be made and probing would fail
		{tHost, -1, probe.Failure, nil},
	}

	prober := New()
	for i, tt := range tests {
		status, _, err := prober.Probe(tt.host, tt.port, 1*time.Second)
		if status != tt.expectedStatus {
			t.Errorf("#%d: expected status=%v, get=%v", i, tt.expectedStatus, status)
		}
		if err != tt.expectedError {
			t.Errorf("#%d: expected error=%v, get=%v", i, tt.expectedError, err)
		}
	}
}
