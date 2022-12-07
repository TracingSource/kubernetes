// +build !providerless
// +build windows

package vsphere

import (
	"fmt"
	"os/exec"
	"strings"
)

func getRawUUID() (string, error) {
	result, err := exec.Command("wmic", "bios", "get", "serialnumber").Output()
	if err != nil {
		return "", err
	}
	lines := strings.FieldsFunc(string(result), func(r rune) bool {
		switch r {
		case '\n', '\r':
			return true
		default:
			return false
		}
	})
	if len(lines) != 2 {
		return "", fmt.Errorf("received unexpected value retrieving vm uuid: %q", string(result))
	}
	return lines[1], nil
}
