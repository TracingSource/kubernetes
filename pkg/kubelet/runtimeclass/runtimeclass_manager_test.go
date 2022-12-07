package runtimeclass_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/kubernetes/pkg/kubelet/runtimeclass"
	rctest "k8s.io/kubernetes/pkg/kubelet/runtimeclass/testing"
	"k8s.io/utils/pointer"
)

func TestLookupRuntimeHandler(t *testing.T) {
	tests := []struct {
		rcn         *string
		expected    string
		expectError bool
	}{
		{rcn: pointer.StringPtr(""), expected: ""},
		{rcn: pointer.StringPtr(rctest.EmptyRuntimeClass), expected: ""},
		{rcn: pointer.StringPtr(rctest.SandboxRuntimeClass), expected: "kata-containers"},
		{rcn: pointer.StringPtr("phantom"), expectError: true},
	}

	manager := runtimeclass.NewManager(rctest.NewPopulatedClient())
	defer rctest.StartManagerSync(manager)()

	for _, test := range tests {
		tname := "nil"
		if test.rcn != nil {
			tname = *test.rcn
		}
		t.Run(fmt.Sprintf("%q->%q(err:%v)", tname, test.expected, test.expectError), func(t *testing.T) {
			handler, err := manager.LookupRuntimeHandler(test.rcn)
			if test.expectError {
				assert.Error(t, err, "handler=%q", handler)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, handler)
			}
		})
	}
}
