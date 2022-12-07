package testing

import (
	"os"

	"k8s.io/kubernetes/pkg/util/sysctl"
)

// fake is a map-backed implementation of sysctl.Interface, for testing/mocking
type fake struct {
	Settings map[string]int
}

func NewFake() *fake {
	return &fake{
		Settings: make(map[string]int),
	}
}

// GetSysctl returns the value for the specified sysctl setting
func (m *fake) GetSysctl(sysctl string) (int, error) {
	v, found := m.Settings[sysctl]
	if !found {
		return -1, os.ErrNotExist
	}
	return v, nil
}

// SetSysctl modifies the specified sysctl flag to the new value
func (m *fake) SetSysctl(sysctl string, newVal int) error {
	m.Settings[sysctl] = newVal
	return nil
}

var _ = sysctl.Interface(&fake{})
