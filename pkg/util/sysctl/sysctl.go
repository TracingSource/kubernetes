package sysctl

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

const (
	sysctlBase = "/proc/sys"
	// VMOvercommitMemory refers to the sysctl variable responsible for defining
	// the memory over-commit policy used by kernel.
	VMOvercommitMemory = "vm/overcommit_memory"
	// VMPanicOnOOM refers to the sysctl variable responsible for defining
	// the OOM behavior used by kernel.
	VMPanicOnOOM = "vm/panic_on_oom"
	// KernelPanic refers to the sysctl variable responsible for defining
	// the timeout after a panic for the kernel to reboot.
	KernelPanic = "kernel/panic"
	// KernelPanicOnOops refers to the sysctl variable responsible for defining
	// the kernel behavior when an oops or BUG is encountered.
	KernelPanicOnOops = "kernel/panic_on_oops"
	// RootMaxKeys refers to the sysctl variable responsible for defining
	// the maximum number of keys that the root user (UID 0 in the root user namespace) may own.
	RootMaxKeys = "kernel/keys/root_maxkeys"
	// RootMaxBytes refers to the sysctl variable responsible for defining
	// the maximum number of bytes of data that the root user (UID 0 in the root user namespace)
	// can hold in the payloads of the keys owned by root.
	RootMaxBytes = "kernel/keys/root_maxbytes"

	// VMOvercommitMemoryAlways represents that kernel performs no memory over-commit handling.
	VMOvercommitMemoryAlways = 1
	// VMPanicOnOOMInvokeOOMKiller represents that kernel calls the oom_killer function when OOM occurs.
	VMPanicOnOOMInvokeOOMKiller = 0

	// KernelPanicOnOopsAlways represents that kernel panics on kernel oops.
	KernelPanicOnOopsAlways = 1
	// KernelPanicRebootTimeout is the timeout seconds after a panic for the kernel to reboot.
	KernelPanicRebootTimeout = 10

	// RootMaxKeysSetting is the maximum number of keys that the root user (UID 0 in the root user namespace) may own.
	// Needed since docker creates a new key per container.
	RootMaxKeysSetting = 1000000
	// RootMaxBytesSetting is the maximum number of bytes of data that the root user (UID 0 in the root user namespace)
	// can hold in the payloads of the keys owned by root.
	// Allocate 25 bytes per key * number of MaxKeys.
	RootMaxBytesSetting = RootMaxKeysSetting * 25
)

// Interface is an injectable interface for running sysctl commands.
type Interface interface {
	// GetSysctl returns the value for the specified sysctl setting
	GetSysctl(sysctl string) (int, error)
	// SetSysctl modifies the specified sysctl flag to the new value
	SetSysctl(sysctl string, newVal int) error
}

// New returns a new Interface for accessing sysctl
func New() Interface {
	return &procSysctl{}
}

// procSysctl implements Interface by reading and writing files under /proc/sys
type procSysctl struct {
}

// GetSysctl returns the value for the specified sysctl setting
func (*procSysctl) GetSysctl(sysctl string) (int, error) {
	data, err := ioutil.ReadFile(path.Join(sysctlBase, sysctl))
	if err != nil {
		return -1, err
	}
	val, err := strconv.Atoi(strings.Trim(string(data), " \n"))
	if err != nil {
		return -1, err
	}
	return val, nil
}

// SetSysctl modifies the specified sysctl flag to the new value
func (*procSysctl) SetSysctl(sysctl string, newVal int) error {
	return ioutil.WriteFile(path.Join(sysctlBase, sysctl), []byte(strconv.Itoa(newVal)), 0640)
}
