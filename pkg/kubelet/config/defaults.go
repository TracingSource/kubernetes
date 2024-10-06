package config

// Defines sane defaults for the kubelet config.
const (
	DefaultKubeletPodsDirName                = "pods" // 该目录下包含当前主机上运行的所有 Pod 列表
	DefaultKubeletVolumesDirName             = "volumes"
	DefaultKubeletVolumeSubpathsDirName      = "volume-subpaths"
	DefaultKubeletVolumeDevicesDirName       = "volumeDevices"
	DefaultKubeletPluginsDirName             = "plugins"
	DefaultKubeletPluginsRegistrationDirName = "plugins_registry"
	DefaultKubeletContainersDirName          = "containers"
	DefaultKubeletPluginContainersDirName    = "plugin-containers"
	DefaultKubeletPodResourcesDirName        = "pod-resources" // 该目录下存在一个 kubelet.sock 文件, 用于与 device plugin 通信
	KubeletPluginsDirSELinuxLabel            = "system_u:object_r:container_file_t:s0"
)
