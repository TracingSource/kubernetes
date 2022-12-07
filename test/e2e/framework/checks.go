package framework

// IsAppArmorSupported checks whether the AppArmor is supported by the node OS distro.
func IsAppArmorSupported() bool {
	return NodeOSDistroIs(AppArmorDistros...)
}
