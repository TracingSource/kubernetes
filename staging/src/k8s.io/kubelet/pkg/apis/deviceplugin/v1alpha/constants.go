package deviceplugin

const (
	// Healthy means that the device is healthy
	Healthy = "Healthy"
	// Unhealthy means that the device is unhealthy
	Unhealthy = "Unhealthy"

	// Version is the current version of the API supported by kubelet
	Version = "v1alpha2"
	// DevicePluginPath is the folder the Device Plugin is expecting sockets to be on
	// Only privileged pods have access to this path
	// Note: Placeholder until we find a "standard path"
	DevicePluginPath = "/var/lib/kubelet/device-plugins/"
	// KubeletSocket is the path of the Kubelet registry socket
	KubeletSocket = DevicePluginPath + "kubelet.sock"
)
