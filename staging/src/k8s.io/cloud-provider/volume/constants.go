package volume

const (
	// ProvisionedVolumeName is the name of a volume in an external cloud
	// that is being provisioned and thus should be ignored by rest of Kubernetes.
	ProvisionedVolumeName = "placeholder-for-provisioning"

	// LabelMultiZoneDelimiter separates zones for volumes
	LabelMultiZoneDelimiter = "__"
)
