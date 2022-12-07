package events

const (
	// volume relevant event reasons
	FailedBinding             = "FailedBinding"
	VolumeMismatch            = "VolumeMismatch"
	VolumeFailedRecycle       = "VolumeFailedRecycle"
	VolumeRecycled            = "VolumeRecycled"
	RecyclerPod               = "RecyclerPod"
	VolumeDelete              = "VolumeDelete"
	VolumeFailedDelete        = "VolumeFailedDelete"
	ExternalProvisioning      = "ExternalProvisioning"
	ProvisioningFailed        = "ProvisioningFailed"
	ProvisioningCleanupFailed = "ProvisioningCleanupFailed"
	ProvisioningSucceeded     = "ProvisioningSucceeded"
	WaitForFirstConsumer      = "WaitForFirstConsumer"
	ExternalExpanding         = "ExternalExpanding"
)
