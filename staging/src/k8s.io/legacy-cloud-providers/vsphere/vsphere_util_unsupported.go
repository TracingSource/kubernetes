// +build !providerless
// +build !windows,!linux

package vsphere

import "fmt"

func getRawUUID() (string, error) {
	return "", fmt.Errorf("Retrieving VM UUID on this build is not implemented.")
}
