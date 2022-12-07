package api

const (
	// TaintExternalCloudProvider sets this taint on a node to mark it as unusable,
	// when kubelet is started with the "external" cloud provider, until a controller
	// from the cloud-controller-manager intitializes this node, and then removes
	// the taint
	TaintExternalCloudProvider = "node.cloudprovider.kubernetes.io/uninitialized"

	// TaintNodeShutdown when node is shutdown in external cloud provider
	TaintNodeShutdown = "node.cloudprovider.kubernetes.io/shutdown"
)
