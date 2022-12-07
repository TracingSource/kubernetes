package azure

import (
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/legacy-cloud-providers/azure"
)

func init() {
	framework.RegisterProvider("azure", newProvider)
}

func newProvider() (framework.ProviderInterface, error) {
	if framework.TestContext.CloudConfig.ConfigFile == "" {
		return nil, fmt.Errorf("config-file must be specified for Azure")
	}
	config, err := os.Open(framework.TestContext.CloudConfig.ConfigFile)
	if err != nil {
		framework.Logf("Couldn't open cloud provider configuration %s: %#v",
			framework.TestContext.CloudConfig.ConfigFile, err)
	}
	defer config.Close()
	azureCloud, err := azure.NewCloud(config)
	return &Provider{
		azureCloud: azureCloud.(*azure.Cloud),
	}, err
}

//Provider is a structure to handle Azure clouds for e2e testing
type Provider struct {
	framework.NullProvider

	azureCloud *azure.Cloud
}

// DeleteNode deletes a node which is specified as the argument
func (p *Provider) DeleteNode(node *v1.Node) error {
	return errors.New("not implemented yet")
}

// CreatePD creates a persistent volume
func (p *Provider) CreatePD(zone string) (string, error) {
	pdName := fmt.Sprintf("%s-%s", framework.TestContext.Prefix, string(uuid.NewUUID()))

	volumeOptions := &azure.ManagedDiskOptions{
		DiskName:           pdName,
		StorageAccountType: compute.StandardLRS,
		ResourceGroup:      "",
		PVCName:            pdName,
		SizeGB:             1,
		Tags:               nil,
		AvailabilityZone:   zone,
		DiskIOPSReadWrite:  "",
		DiskMBpsReadWrite:  "",
	}
	return p.azureCloud.CreateManagedDisk(volumeOptions)
}

// DeletePD deletes a persistent volume
func (p *Provider) DeletePD(pdName string) error {
	if err := p.azureCloud.DeleteManagedDisk(pdName); err != nil {
		framework.Logf("failed to delete Azure volume %q: %v", pdName, err)
		return err
	}
	return nil
}

// EnableAndDisableInternalLB returns functions for both enabling and disabling internal Load Balancer
func (p *Provider) EnableAndDisableInternalLB() (enable, disable func(svc *v1.Service)) {
	enable = func(svc *v1.Service) {
		svc.ObjectMeta.Annotations = map[string]string{azure.ServiceAnnotationLoadBalancerInternal: "true"}
	}
	disable = func(svc *v1.Service) {
		svc.ObjectMeta.Annotations = map[string]string{azure.ServiceAnnotationLoadBalancerInternal: "false"}
	}
	return
}
