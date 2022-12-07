// Get information about the cloud provider (if any) cAdvisor is running on.

package cloudinfo

import (
	info "github.com/google/cadvisor/info/v1"
	"k8s.io/klog"
)

type CloudInfo interface {
	GetCloudProvider() info.CloudProvider
	GetInstanceType() info.InstanceType
	GetInstanceID() info.InstanceID
}

// CloudProvider is an abstraction for providing cloud-specific information.
type CloudProvider interface {
	// IsActiveProvider determines whether this is the cloud provider operating
	// this instance.
	IsActiveProvider() bool
	// GetInstanceType gets the type of instance this process is running on.
	// The behavior is undefined if this is not the active provider.
	GetInstanceType() info.InstanceType
	// GetInstanceType gets the ID of the instance this process is running on.
	// The behavior is undefined if this is not the active provider.
	GetInstanceID() info.InstanceID
}

var providers = map[info.CloudProvider]CloudProvider{}

// RegisterCloudProvider registers the given cloud provider
func RegisterCloudProvider(name info.CloudProvider, provider CloudProvider) {
	if _, alreadyRegistered := providers[name]; alreadyRegistered {
		klog.Warningf("Duplicate registration of CloudProvider %s", name)
	}
	providers[name] = provider
}

type realCloudInfo struct {
	cloudProvider info.CloudProvider
	instanceType  info.InstanceType
	instanceID    info.InstanceID
}

func NewRealCloudInfo() CloudInfo {
	for name, provider := range providers {
		if provider.IsActiveProvider() {
			return &realCloudInfo{
				cloudProvider: name,
				instanceType:  provider.GetInstanceType(),
				instanceID:    provider.GetInstanceID(),
			}
		}
	}

	// No registered active provider.
	return &realCloudInfo{
		cloudProvider: info.UnknownProvider,
		instanceType:  info.UnknownInstance,
		instanceID:    info.UnNamedInstance,
	}
}

func (self *realCloudInfo) GetCloudProvider() info.CloudProvider {
	return self.cloudProvider
}

func (self *realCloudInfo) GetInstanceType() info.InstanceType {
	return self.instanceType
}

func (self *realCloudInfo) GetInstanceID() info.InstanceID {
	return self.instanceID
}
