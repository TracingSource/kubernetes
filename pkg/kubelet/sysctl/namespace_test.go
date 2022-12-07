package sysctl

import (
	"testing"
)

func TestNamespacedBy(t *testing.T) {
	tests := map[string]Namespace{
		"kernel.shm_rmid_forced": IpcNamespace,
		"net.a.b.c":              NetNamespace,
		"fs.mqueue.a.b.c":        IpcNamespace,
		"foo":                    UnknownNamespace,
	}

	for sysctl, ns := range tests {
		if got := NamespacedBy(sysctl); got != ns {
			t.Errorf("wrong namespace for %q: got=%s want=%s", sysctl, got, ns)
		}
	}
}
