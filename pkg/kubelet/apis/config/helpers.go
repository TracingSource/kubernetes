package config

// KubeletConfigurationPathRefs returns pointers to all of the KubeletConfiguration fields that contain filepaths.
// You might use this, for example, to resolve all relative paths against some common root before
// passing the configuration to the application. This method must be kept up to date as new fields are added.
func KubeletConfigurationPathRefs(kc *KubeletConfiguration) []*string {
	paths := []*string{}
	paths = append(paths, &kc.StaticPodPath)
	paths = append(paths, &kc.Authentication.X509.ClientCAFile)
	paths = append(paths, &kc.TLSCertFile)
	paths = append(paths, &kc.TLSPrivateKeyFile)
	paths = append(paths, &kc.ResolverConfig)
	return paths
}
