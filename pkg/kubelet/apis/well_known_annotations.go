package apis

const (
	// AnnotationProvidedIPAddr is a node IP annotation set by the "external" cloud provider.
	// When kubelet is started with the "external" cloud provider, then
	// it sets this annotation on the node to denote an ip address set from the
	// cmd line flag (--node-ip). This ip is verified with the cloudprovider as valid by
	// the cloud-controller-manager
	AnnotationProvidedIPAddr = "alpha.kubernetes.io/provided-node-ip"
)
