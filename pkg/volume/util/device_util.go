package util

//DeviceUtil is a util for common device methods
type DeviceUtil interface {
	FindMultipathDeviceForDevice(disk string) string
	FindSlaveDevicesOnMultipath(disk string) []string
	GetISCSIPortalHostMapForTarget(targetIqn string) (map[string]int, error)
	FindDevicesForISCSILun(targetIqn string, lun int) ([]string, error)
}

type deviceHandler struct {
	getIo IoUtil
}

//NewDeviceHandler Create a new IoHandler implementation
func NewDeviceHandler(io IoUtil) DeviceUtil {
	return &deviceHandler{getIo: io}
}
