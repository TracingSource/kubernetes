package cloudinfo

import (
	"io/ioutil"
	"strings"

	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/utils/cloudinfo"
)

const (
	sysVendorFileName    = "/sys/class/dmi/id/sys_vendor"
	biosUUIDFileName     = "/sys/class/dmi/id/product_uuid"
	microsoftCorporation = "Microsoft Corporation"
)

func init() {
	cloudinfo.RegisterCloudProvider(info.Azure, &provider{})
}

type provider struct{}

var _ cloudinfo.CloudProvider = provider{}

func (provider) IsActiveProvider() bool {
	data, err := ioutil.ReadFile(sysVendorFileName)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), microsoftCorporation)
}

// TODO: Implement method.
func (provider) GetInstanceType() info.InstanceType {
	return info.UnknownInstance
}

func (provider) GetInstanceID() info.InstanceID {
	data, err := ioutil.ReadFile(biosUUIDFileName)
	if err != nil {
		return info.UnNamedInstance
	}
	return info.InstanceID(strings.TrimSuffix(string(data), "\n"))
}
