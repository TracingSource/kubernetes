package vclib

import (
	"time"

	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

// Cloud Provider API constants
const (
	APICreateVolume = "CreateVolume"
	APIDeleteVolume = "DeleteVolume"
	APIAttachVolume = "AttachVolume"
	APIDetachVolume = "DetachVolume"
)

// Cloud Provider Operation constants
const (
	OperationDeleteVolume                  = "DeleteVolumeOperation"
	OperationAttachVolume                  = "AttachVolumeOperation"
	OperationDetachVolume                  = "DetachVolumeOperation"
	OperationDiskIsAttached                = "DiskIsAttachedOperation"
	OperationDisksAreAttached              = "DisksAreAttachedOperation"
	OperationCreateVolume                  = "CreateVolumeOperation"
	OperationCreateVolumeWithPolicy        = "CreateVolumeWithPolicyOperation"
	OperationCreateVolumeWithRawVSANPolicy = "CreateVolumeWithRawVSANPolicyOperation"
)

// vsphereAPIMetric is for recording latency of Single API Call.
var vsphereAPIMetric = metrics.NewHistogramVec(
	&metrics.HistogramOpts{
		Name:           "cloudprovider_vsphere_api_request_duration_seconds",
		Help:           "Latency of vsphere api call",
		StabilityLevel: metrics.ALPHA,
	},
	[]string{"request"},
)

var vsphereAPIErrorMetric = metrics.NewCounterVec(
	&metrics.CounterOpts{
		Name:           "cloudprovider_vsphere_api_request_errors",
		Help:           "vsphere Api errors",
		StabilityLevel: metrics.ALPHA,
	},
	[]string{"request"},
)

// vsphereOperationMetric is for recording latency of vSphere Operation which invokes multiple APIs to get the task done.
var vsphereOperationMetric = metrics.NewHistogramVec(
	&metrics.HistogramOpts{
		Name:           "cloudprovider_vsphere_operation_duration_seconds",
		Help:           "Latency of vsphere operation call",
		StabilityLevel: metrics.ALPHA,
	},
	[]string{"operation"},
)

var vsphereOperationErrorMetric = metrics.NewCounterVec(
	&metrics.CounterOpts{
		Name:           "cloudprovider_vsphere_operation_errors",
		Help:           "vsphere operation errors",
		StabilityLevel: metrics.ALPHA,
	},
	[]string{"operation"},
)

// RegisterMetrics registers all the API and Operation metrics
func RegisterMetrics() {
	legacyregistry.MustRegister(vsphereAPIMetric)
	legacyregistry.MustRegister(vsphereAPIErrorMetric)
	legacyregistry.MustRegister(vsphereOperationMetric)
	legacyregistry.MustRegister(vsphereOperationErrorMetric)
}

// RecordvSphereMetric records the vSphere API and Operation metrics
func RecordvSphereMetric(actionName string, requestTime time.Time, err error) {
	switch actionName {
	case APICreateVolume, APIDeleteVolume, APIAttachVolume, APIDetachVolume:
		recordvSphereAPIMetric(actionName, requestTime, err)
	default:
		recordvSphereOperationMetric(actionName, requestTime, err)
	}
}

func recordvSphereAPIMetric(actionName string, requestTime time.Time, err error) {
	if err != nil {
		vsphereAPIErrorMetric.With(metrics.Labels{"request": actionName}).Inc()
	} else {
		vsphereAPIMetric.With(metrics.Labels{"request": actionName}).Observe(calculateTimeTaken(requestTime))
	}
}

func recordvSphereOperationMetric(actionName string, requestTime time.Time, err error) {
	if err != nil {
		vsphereOperationErrorMetric.With(metrics.Labels{"operation": actionName}).Inc()
	} else {
		vsphereOperationMetric.With(metrics.Labels{"operation": actionName}).Observe(calculateTimeTaken(requestTime))
	}
}

// RecordCreateVolumeMetric records the Create Volume metric
func RecordCreateVolumeMetric(volumeOptions *VolumeOptions, requestTime time.Time, err error) {
	var actionName string
	if volumeOptions.StoragePolicyName != "" {
		actionName = OperationCreateVolumeWithPolicy
	} else if volumeOptions.VSANStorageProfileData != "" {
		actionName = OperationCreateVolumeWithRawVSANPolicy
	} else {
		actionName = OperationCreateVolume
	}
	RecordvSphereMetric(actionName, requestTime, err)
}

func calculateTimeTaken(requestBeginTime time.Time) (timeTaken float64) {
	if !requestBeginTime.IsZero() {
		timeTaken = time.Since(requestBeginTime).Seconds()
	} else {
		timeTaken = 0
	}
	return timeTaken
}
