package kuberuntime

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"k8s.io/component-base/metrics/legacyregistry"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
)

func TestRecordOperation(t *testing.T) {
	legacyregistry.MustRegister(metrics.RuntimeOperations)
	legacyregistry.MustRegister(metrics.RuntimeOperationsDuration)
	legacyregistry.MustRegister(metrics.RuntimeOperationsErrors)

	temporalServer := "127.0.0.1:1234"
	l, err := net.Listen("tcp", temporalServer)
	assert.NoError(t, err)
	defer l.Close()

	prometheusURL := "http://" + temporalServer + "/metrics"
	mux := http.NewServeMux()
	//lint:ignore SA1019 ignore deprecated warning until we move off of global registries
	mux.Handle("/metrics", legacyregistry.Handler())
	server := &http.Server{
		Addr:    temporalServer,
		Handler: mux,
	}
	go func() {
		server.Serve(l)
	}()

	recordOperation("create_container", time.Now())
	runtimeOperationsCounterExpected := "kubelet_runtime_operations_total{operation_type=\"create_container\"} 1"
	runtimeOperationsDurationExpected := "kubelet_runtime_operations_duration_seconds_count{operation_type=\"create_container\"} 1"

	assert.HTTPBodyContains(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}), "GET", prometheusURL, nil, runtimeOperationsCounterExpected)

	assert.HTTPBodyContains(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}), "GET", prometheusURL, nil, runtimeOperationsDurationExpected)
}

func TestInstrumentedVersion(t *testing.T) {
	fakeRuntime, _, _, _ := createTestRuntimeManager()
	irs := newInstrumentedRuntimeService(fakeRuntime)
	vr, err := irs.Version("1")
	assert.NoError(t, err)
	assert.Equal(t, kubeRuntimeAPIVersion, vr.Version)
}

func TestStatus(t *testing.T) {
	fakeRuntime, _, _, _ := createTestRuntimeManager()
	fakeRuntime.FakeStatus = &runtimeapi.RuntimeStatus{
		Conditions: []*runtimeapi.RuntimeCondition{
			{Type: runtimeapi.RuntimeReady, Status: false},
			{Type: runtimeapi.NetworkReady, Status: true},
		},
	}
	irs := newInstrumentedRuntimeService(fakeRuntime)
	actural, err := irs.Status()
	assert.NoError(t, err)
	expected := &runtimeapi.RuntimeStatus{
		Conditions: []*runtimeapi.RuntimeCondition{
			{Type: runtimeapi.RuntimeReady, Status: false},
			{Type: runtimeapi.NetworkReady, Status: true},
		},
	}
	assert.Equal(t, expected, actural)
}
