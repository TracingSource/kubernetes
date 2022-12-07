package sysctl

import (
	"strings"
)

// Namespace represents a kernel namespace name.
type Namespace string

const (
	// the Linux IPC namespace
	IpcNamespace = Namespace("ipc")

	// the network namespace
	NetNamespace = Namespace("net")

	// the zero value if no namespace is known
	UnknownNamespace = Namespace("")
)

var namespaces = map[string]Namespace{
	"kernel.sem": IpcNamespace,
}

var prefixNamespaces = map[string]Namespace{
	"kernel.shm": IpcNamespace,
	"kernel.msg": IpcNamespace,
	"fs.mqueue.": IpcNamespace,
	"net.":       NetNamespace,
}

// NamespacedBy returns the namespace of the Linux kernel for a sysctl, or
// UnknownNamespace if the sysctl is not known to be namespaced.
func NamespacedBy(val string) Namespace {
	if ns, found := namespaces[val]; found {
		return ns
	}
	for p, ns := range prefixNamespaces {
		if strings.HasPrefix(val, p) {
			return ns
		}
	}
	return UnknownNamespace
}
