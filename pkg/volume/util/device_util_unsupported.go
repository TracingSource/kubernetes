// +build !linux

package util

// FindMultipathDeviceForDevice unsupported returns ""
func (handler *deviceHandler) FindMultipathDeviceForDevice(device string) string {
	return ""
}

// FindSlaveDevicesOnMultipath unsupported returns ""
func (handler *deviceHandler) FindSlaveDevicesOnMultipath(disk string) []string {
	out := []string{}
	return out
}

// GetISCSIPortalHostMapForTarget unsupported returns nil
func (handler *deviceHandler) GetISCSIPortalHostMapForTarget(targetIqn string) (map[string]int, error) {
	portalHostMap := make(map[string]int)
	return portalHostMap, nil
}

// FindDevicesForISCSILun unsupported returns nil
func (handler *deviceHandler) FindDevicesForISCSILun(targetIqn string, lun int) ([]string, error) {
	devices := []string{}
	return devices, nil
}
