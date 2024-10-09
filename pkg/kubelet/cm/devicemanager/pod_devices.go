package devicemanager

import (
	"k8s.io/klog"

	"k8s.io/apimachinery/pkg/util/sets"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	podresourcesapi "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/cm/devicemanager/checkpoint"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

type deviceAllocateInfo struct {
	// deviceIds contains device Ids allocated to this container for the given resourceName.
	deviceIds sets.String
	// allocResp contains cached rpc AllocateResponse.
	allocResp *pluginapi.ContainerAllocateResponse
}

type resourceAllocateInfo map[string]deviceAllocateInfo // Keyed by resourceName.
type containerDevices map[string]resourceAllocateInfo   // Keyed by containerName.
// podDevices 中存储着所有已分配的 devices 信息, 并按 pod/container 进行分级归类,
// 作为 kubelet 的本地缓存.
type podDevices map[string]containerDevices             // Keyed by podUID.

func (pdev podDevices) pods() sets.String {
	ret := sets.NewString()
	for k := range pdev {
		ret.Insert(k)
	}
	return ret
}

func (pdev podDevices) insert(
	podUID, contName, resource string, devices sets.String, 
	resp *pluginapi.ContainerAllocateResponse,
) {
	if _, podExists := pdev[podUID]; !podExists {
		pdev[podUID] = make(containerDevices)
	}
	if _, contExists := pdev[podUID][contName]; !contExists {
		pdev[podUID][contName] = make(resourceAllocateInfo)
	}
	pdev[podUID][contName][resource] = deviceAllocateInfo{
		deviceIds: devices,
		allocResp: resp,
	}
}

func (pdev podDevices) delete(pods []string) {
	for _, uid := range pods {
		delete(pdev, uid)
	}
}

// containerDevices 从缓存中获取已分配给该 pod/container 下的设备列表并返回.
// 容器重启后可能会尝试重新分配, 如果这里能查到, 到主调函数就不需要继续往下走了.
//
// Returns list of device Ids allocated to the given container for the given resource.
// Returns nil if we don't have cached state for the given <podUID, contName, resource>.
func (pdev podDevices) containerDevices(podUID, contName, resource string) sets.String {
	if _, podExists := pdev[podUID]; !podExists {
		return nil
	}
	if _, contExists := pdev[podUID][contName]; !contExists {
		return nil
	}
	devs, resourceExists := pdev[podUID][contName][resource]
	if !resourceExists {
		return nil
	}
	return devs.deviceIds
}

// Populates allocatedResources with the device resources allocated to the specified <podUID, contName>.
func (pdev podDevices) addContainerAllocatedResources(
	podUID, contName string, allocatedResources map[string]sets.String,
) {
	containers, exists := pdev[podUID]
	if !exists {
		return
	}
	resources, exists := containers[contName]
	if !exists {
		return
	}
	for resource, devices := range resources {
		allocatedResources[resource] = allocatedResources[resource].Union(devices.deviceIds)
	}
}

// Removes the device resources allocated to the specified <podUID, contName> from allocatedResources.
func (pdev podDevices) removeContainerAllocatedResources(
	podUID, contName string, allocatedResources map[string]sets.String,
) {
	containers, exists := pdev[podUID]
	if !exists {
		return
	}
	resources, exists := containers[contName]
	if !exists {
		return
	}
	for resource, devices := range resources {
		allocatedResources[resource] = allocatedResources[resource].Difference(devices.deviceIds)
	}
}

// devices 从 podDevices 中取出所有的设备信息, 作为集合返回.
// (这些设备都是已分配的设备, 属于某个 Pod)
//
// 	@return ret: key 为扩展资源名称(与 cpu/memory 同级), value 为 podDevices 中
// 所有的设备的 ID 集合.
//
// Returns all of devices allocated to the pods being tracked, keyed by resourceName.
func (pdev podDevices) devices() map[string]sets.String {
	ret := make(map[string]sets.String)
	for _, containerDevices := range pdev {
		for _, resources := range containerDevices {
			for resource, devices := range resources {
				if _, exists := ret[resource]; !exists {
					ret[resource] = sets.NewString()
				}
				if devices.allocResp != nil {
					ret[resource] = ret[resource].Union(devices.deviceIds)
				}
			}
		}
	}
	return ret
}

// Turns podDevices to checkpointData.
func (pdev podDevices) toCheckpointData() []checkpoint.PodDevicesEntry {
	var data []checkpoint.PodDevicesEntry
	for podUID, containerDevices := range pdev {
		for conName, resources := range containerDevices {
			for resource, devices := range resources {
				devIds := devices.deviceIds.UnsortedList()
				if devices.allocResp == nil {
					klog.Errorf(
						"Can't marshal allocResp for %v %v %v: allocation response is missing", 
						podUID, conName, resource,
					)
					continue
				}

				allocResp, err := devices.allocResp.Marshal()
				if err != nil {
					klog.Errorf(
						"Can't marshal allocResp for %v %v %v: %v", 
						podUID, conName, resource, err,
					)
					continue
				}
				data = append(data, checkpoint.PodDevicesEntry{
					PodUID:        podUID,
					ContainerName: conName,
					ResourceName:  resource,
					DeviceIDs:     devIds,
					AllocResp:     allocResp})
			}
		}
	}
	return data
}

// Populates podDevices from the passed in checkpointData.
func (pdev podDevices) fromCheckpointData(data []checkpoint.PodDevicesEntry) {
	for _, entry := range data {
		klog.V(2).Infof(
			"Get checkpoint entry: %v %v %v %v %v\n",
			entry.PodUID, entry.ContainerName, 
			entry.ResourceName, entry.DeviceIDs, entry.AllocResp,
		)
		devIDs := sets.NewString()
		for _, devID := range entry.DeviceIDs {
			devIDs.Insert(devID)
		}
		allocResp := &pluginapi.ContainerAllocateResponse{}
		err := allocResp.Unmarshal(entry.AllocResp)
		if err != nil {
			klog.Errorf(
				"Can't unmarshal allocResp for %v %v %v: %v", 
				entry.PodUID, entry.ContainerName, entry.ResourceName, err,
			)
			continue
		}
		pdev.insert(entry.PodUID, entry.ContainerName, entry.ResourceName, devIDs, allocResp)
	}
}

// Returns combined container runtime settings to consume the container's allocated devices.
func (pdev podDevices) deviceRunContainerOptions(
	podUID, contName string,
) *DeviceRunContainerOptions {
	containers, exists := pdev[podUID]
	if !exists {
		return nil
	}
	resources, exists := containers[contName]
	if !exists {
		return nil
	}
	opts := &DeviceRunContainerOptions{}
	// Maps to detect duplicate settings.
	devsMap := make(map[string]string)
	mountsMap := make(map[string]string)
	envsMap := make(map[string]string)
	annotationsMap := make(map[string]string)
	// Loops through AllocationResponses of all cached device resources.
	for _, devices := range resources {
		resp := devices.allocResp
		// Each Allocate response has the following artifacts.
		// Environment variables
		// Mount points
		// Device files
		// Container annotations
		// These artifacts are per resource per container.
		// Updates RunContainerOptions.Envs.
		for k, v := range resp.Envs {
			if e, ok := envsMap[k]; ok {
				klog.V(4).Infof("Skip existing env %s %s", k, v)
				if e != v {
					klog.Errorf("Environment variable %s has conflicting setting: %s and %s", k, e, v)
				}
				continue
			}
			klog.V(4).Infof("Add env %s %s", k, v)
			envsMap[k] = v
			opts.Envs = append(opts.Envs, kubecontainer.EnvVar{Name: k, Value: v})
		}

		// Updates RunContainerOptions.Devices.
		for _, dev := range resp.Devices {
			if d, ok := devsMap[dev.ContainerPath]; ok {
				klog.V(4).Infof("Skip existing device %s %s", dev.ContainerPath, dev.HostPath)
				if d != dev.HostPath {
					klog.Errorf(
						"Container device %s has conflicting mapping host devices: %s and %s",
						dev.ContainerPath, d, dev.HostPath,
					)
				}
				continue
			}
			klog.V(4).Infof("Add device %s %s", dev.ContainerPath, dev.HostPath)
			devsMap[dev.ContainerPath] = dev.HostPath
			opts.Devices = append(opts.Devices, kubecontainer.DeviceInfo{
				PathOnHost:      dev.HostPath,
				PathInContainer: dev.ContainerPath,
				Permissions:     dev.Permissions,
			})
		}

		// Updates RunContainerOptions.Mounts.
		for _, mount := range resp.Mounts {
			if m, ok := mountsMap[mount.ContainerPath]; ok {
				klog.V(4).Infof("Skip existing mount %s %s", mount.ContainerPath, mount.HostPath)
				if m != mount.HostPath {
					klog.Errorf("Container mount %s has conflicting mapping host mounts: %s and %s",
						mount.ContainerPath, m, mount.HostPath)
				}
				continue
			}
			klog.V(4).Infof("Add mount %s %s", mount.ContainerPath, mount.HostPath)
			mountsMap[mount.ContainerPath] = mount.HostPath
			opts.Mounts = append(opts.Mounts, kubecontainer.Mount{
				Name:          mount.ContainerPath,
				ContainerPath: mount.ContainerPath,
				HostPath:      mount.HostPath,
				ReadOnly:      mount.ReadOnly,
				// TODO: This may need to be part of Device plugin API.
				SELinuxRelabel: false,
			})
		}

		// Updates for Annotations
		for k, v := range resp.Annotations {
			if e, ok := annotationsMap[k]; ok {
				klog.V(4).Infof("Skip existing annotation %s %s", k, v)
				if e != v {
					klog.Errorf("Annotation %s has conflicting setting: %s and %s", k, e, v)
				}
				continue
			}
			klog.V(4).Infof("Add annotation %s %s", k, v)
			annotationsMap[k] = v
			opts.Annotations = append(opts.Annotations, kubecontainer.Annotation{Name: k, Value: v})
		}
	}
	return opts
}

// getContainerDevices returns the devices assigned to the provided container for all ResourceNames
func (pdev podDevices) getContainerDevices(podUID, contName string) []*podresourcesapi.ContainerDevices {
	if _, podExists := pdev[podUID]; !podExists {
		return nil
	}
	if _, contExists := pdev[podUID][contName]; !contExists {
		return nil
	}
	cDev := []*podresourcesapi.ContainerDevices{}
	for resource, allocateInfo := range pdev[podUID][contName] {
		cDev = append(cDev, &podresourcesapi.ContainerDevices{
			ResourceName: resource,
			DeviceIds:    allocateInfo.deviceIds.UnsortedList(),
		})
	}
	return cDev
}
