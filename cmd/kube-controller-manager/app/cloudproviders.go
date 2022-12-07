package app

import (
	"fmt"

	"k8s.io/klog"

	"k8s.io/client-go/informers"
	cloudprovider "k8s.io/cloud-provider"
)

// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> CreateControllerContext()
//
// createCloudProvider helps consolidate what is needed for cloud providers, 
// we explicitly list the things that the cloud providers need as parameters,
// so we can control
func createCloudProvider(
	cloudProvider string, externalCloudVolumePlugin string, cloudConfigFile string,
	allowUntaggedCloud bool, sharedInformers informers.SharedInformerFactory,
) (cloudprovider.Interface, ControllerLoopMode, error) {
	var cloud cloudprovider.Interface
	var loopMode ControllerLoopMode
	var err error
	if cloudprovider.IsExternal(cloudProvider) {
		loopMode = ExternalLoops
		if externalCloudVolumePlugin == "" {
			// externalCloudVolumePlugin is temporary until we split all cloud providers out.
			// So we just tell the caller that we need to run ExternalLoops without any cloud provider.
			return nil, loopMode, nil
		}
		cloud, err = cloudprovider.InitCloudProvider(externalCloudVolumePlugin, cloudConfigFile)
	} else {
		loopMode = IncludeCloudLoops
		cloud, err = cloudprovider.InitCloudProvider(cloudProvider, cloudConfigFile)
	}
	if err != nil {
		return nil, loopMode, fmt.Errorf("cloud provider could not be initialized: %v", err)
	}

	if cloud != nil && !cloud.HasClusterID() {
		if allowUntaggedCloud {
			klog.Warning("detected a cluster without a ClusterID.  A ClusterID will be required in the future.  Please tag your cluster to avoid any future issues")
		} else {
			return nil, loopMode, fmt.Errorf("no ClusterID Found.  A ClusterID is required for the cloud provider to function properly.  This check can be bypassed by setting the allow-untagged-cloud option")
		}
	}

	if informerUserCloud, ok := cloud.(cloudprovider.InformerUser); ok {
		informerUserCloud.SetInformers(sharedInformers)
	}
	return cloud, loopMode, err
}
