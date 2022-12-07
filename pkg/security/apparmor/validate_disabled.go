// +build !linux

package apparmor

func init() {
	// If Kubernetes was not built for linux, apparmor is always disabled.
	isDisabledBuild = true
}
