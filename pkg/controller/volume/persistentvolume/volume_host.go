package persistentvolume

import (
	"fmt"
	"net"

	"k8s.io/klog"
	utilexec "k8s.io/utils/exec"
	"k8s.io/utils/mount"

	authenticationv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	cloudprovider "k8s.io/cloud-provider"
	vol "k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util/subpath"
)

// VolumeHost interface implementation for PersistentVolumeController.

var _ vol.VolumeHost = &PersistentVolumeController{}

func (ctrl *PersistentVolumeController) GetPluginDir(pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetVolumeDevicePluginDir(pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodsDir() string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodVolumeDir(podUID types.UID, pluginName string, volumeName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodPluginDir(podUID types.UID, pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetPodVolumeDeviceDir(ppodUID types.UID, pluginName string) string {
	return ""
}

func (ctrl *PersistentVolumeController) GetKubeClient() clientset.Interface {
	return ctrl.kubeClient
}

func (ctrl *PersistentVolumeController) NewWrapperMounter(volName string, spec vol.Spec, pod *v1.Pod, opts vol.VolumeOptions) (vol.Mounter, error) {
	return nil, fmt.Errorf("PersistentVolumeController.NewWrapperMounter is not implemented")
}

func (ctrl *PersistentVolumeController) NewWrapperUnmounter(volName string, spec vol.Spec, podUID types.UID) (vol.Unmounter, error) {
	return nil, fmt.Errorf("PersistentVolumeController.NewWrapperMounter is not implemented")
}

func (ctrl *PersistentVolumeController) GetCloudProvider() cloudprovider.Interface {
	return ctrl.cloud
}

func (ctrl *PersistentVolumeController) GetMounter(pluginName string) mount.Interface {
	return nil
}

func (ctrl *PersistentVolumeController) GetHostName() string {
	return ""
}

func (ctrl *PersistentVolumeController) GetHostIP() (net.IP, error) {
	return nil, fmt.Errorf("PersistentVolumeController.GetHostIP() is not implemented")
}

func (ctrl *PersistentVolumeController) GetNodeAllocatable() (v1.ResourceList, error) {
	return v1.ResourceList{}, nil
}

func (ctrl *PersistentVolumeController) GetSecretFunc() func(namespace, name string) (*v1.Secret, error) {
	return func(_, _ string) (*v1.Secret, error) {
		return nil, fmt.Errorf("GetSecret unsupported in PersistentVolumeController")
	}
}

func (ctrl *PersistentVolumeController) GetConfigMapFunc() func(namespace, name string) (*v1.ConfigMap, error) {
	return func(_, _ string) (*v1.ConfigMap, error) {
		return nil, fmt.Errorf("GetConfigMap unsupported in PersistentVolumeController")
	}
}

func (ctrl *PersistentVolumeController) GetServiceAccountTokenFunc() func(_, _ string, _ *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
	return func(_, _ string, _ *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
		return nil, fmt.Errorf("GetServiceAccountToken unsupported in PersistentVolumeController")
	}
}

func (ctrl *PersistentVolumeController) DeleteServiceAccountTokenFunc() func(types.UID) {
	return func(types.UID) {
		klog.Errorf("DeleteServiceAccountToken unsupported in PersistentVolumeController")
	}
}

func (adc *PersistentVolumeController) GetExec(pluginName string) utilexec.Interface {
	return utilexec.New()
}

func (ctrl *PersistentVolumeController) GetNodeLabels() (map[string]string, error) {
	return nil, fmt.Errorf("GetNodeLabels() unsupported in PersistentVolumeController")
}

func (ctrl *PersistentVolumeController) GetNodeName() types.NodeName {
	return ""
}

func (ctrl *PersistentVolumeController) GetEventRecorder() record.EventRecorder {
	return ctrl.eventRecorder
}

func (ctrl *PersistentVolumeController) GetSubpather() subpath.Interface {
	// No volume plugin needs Subpaths in PV controller.
	return nil
}
