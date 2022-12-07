package util

const (
	// PVCProtectionFinalizer is the name of finalizer on PVCs that have a running pod.
	PVCProtectionFinalizer = "kubernetes.io/pvc-protection"

	// PVProtectionFinalizer is the name of finalizer on PVs that are bound by PVCs
	PVProtectionFinalizer = "kubernetes.io/pv-protection"
)
