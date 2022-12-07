// +build !providerless

package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// NewNetworkV2 creates a ServiceClient that may be used with the neutron v2 API
func (os *OpenStack) NewNetworkV2() (*gophercloud.ServiceClient, error) {
	network, err := openstack.NewNetworkV2(os.provider, gophercloud.EndpointOpts{
		Region: os.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find network v2 endpoint for region %s: %v", os.region, err)
	}
	return network, nil
}

// NewComputeV2 creates a ServiceClient that may be used with the nova v2 API
func (os *OpenStack) NewComputeV2() (*gophercloud.ServiceClient, error) {
	compute, err := openstack.NewComputeV2(os.provider, gophercloud.EndpointOpts{
		Region: os.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find compute v2 endpoint for region %s: %v", os.region, err)
	}
	return compute, nil
}

// NewBlockStorageV1 creates a ServiceClient that may be used with the Cinder v1 API
func (os *OpenStack) NewBlockStorageV1() (*gophercloud.ServiceClient, error) {
	storage, err := openstack.NewBlockStorageV1(os.provider, gophercloud.EndpointOpts{
		Region: os.region,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cinder v1 client for region %s: %v", os.region, err)
	}
	return storage, nil
}

// NewBlockStorageV2 creates a ServiceClient that may be used with the Cinder v2 API
func (os *OpenStack) NewBlockStorageV2() (*gophercloud.ServiceClient, error) {
	storage, err := openstack.NewBlockStorageV2(os.provider, gophercloud.EndpointOpts{
		Region: os.region,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cinder v2 client for region %s: %v", os.region, err)
	}
	return storage, nil
}

// NewBlockStorageV3 creates a ServiceClient that may be used with the Cinder v3 API
func (os *OpenStack) NewBlockStorageV3() (*gophercloud.ServiceClient, error) {
	storage, err := openstack.NewBlockStorageV3(os.provider, gophercloud.EndpointOpts{
		Region: os.region,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cinder v3 client for region %s: %v", os.region, err)
	}
	return storage, nil
}

// NewLoadBalancerV2 creates a ServiceClient that may be used with the Neutron LBaaS v2 API
func (os *OpenStack) NewLoadBalancerV2() (*gophercloud.ServiceClient, error) {
	var lb *gophercloud.ServiceClient
	var err error
	if os.lbOpts.UseOctavia {
		lb, err = openstack.NewLoadBalancerV2(os.provider, gophercloud.EndpointOpts{
			Region: os.region,
		})
	} else {
		lb, err = openstack.NewNetworkV2(os.provider, gophercloud.EndpointOpts{
			Region: os.region,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find load-balancer v2 endpoint for region %s: %v", os.region, err)
	}
	return lb, nil
}
