package kubelet

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/api/v1/resource"
	podshelper "k8s.io/kubernetes/pkg/apis/core/pods"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/fieldpath"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/envvars"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	"k8s.io/kubernetes/pkg/kubelet/images"
	"k8s.io/kubernetes/pkg/kubelet/server/portforward"
	remotecommandserver "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
	"k8s.io/kubernetes/pkg/kubelet/status"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/kubelet/util/format"
	volumeutil "k8s.io/kubernetes/pkg/volume/util"
	"k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
	"k8s.io/kubernetes/third_party/forked/golang/expansion"
)

const (
	managedHostsHeader                = "# Kubernetes-managed hosts file.\n"
	managedHostsHeaderWithHostNetwork = "# Kubernetes-managed hosts file (host network).\n"
)

// Get a list of pods that have data directories.
func (kl *Kubelet) listPodsFromDisk() ([]types.UID, error) {
	podInfos, err := ioutil.ReadDir(kl.getPodsDir())
	if err != nil {
		return nil, err
	}
	pods := []types.UID{}
	for i := range podInfos {
		if podInfos[i].IsDir() {
			pods = append(pods, types.UID(podInfos[i].Name()))
		}
	}
	return pods, nil
}

// GetActivePods 返回非 terminating 状态的 pod 列表.
//
// GetActivePods returns non-terminal pods
func (kl *Kubelet) GetActivePods() []*v1.Pod {
	allPods := kl.podManager.GetPods()
	activePods := kl.filterOutTerminatedPods(allPods)
	return activePods
}

// makeBlockVolumes maps the raw block devices specified in the path of the container
// Experimental
func (kl *Kubelet) makeBlockVolumes(pod *v1.Pod, container *v1.Container, podVolumes kubecontainer.VolumeMap, blkutil volumepathhandler.BlockVolumePathHandler) ([]kubecontainer.DeviceInfo, error) {
	var devices []kubecontainer.DeviceInfo
	for _, device := range container.VolumeDevices {
		// check path is absolute
		if !filepath.IsAbs(device.DevicePath) {
			return nil, fmt.Errorf("error DevicePath `%s` must be an absolute path", device.DevicePath)
		}
		vol, ok := podVolumes[device.Name]
		if !ok || vol.BlockVolumeMapper == nil {
			klog.Errorf("Block volume cannot be satisfied for container %q, because the volume is missing or the volume mapper is nil: %+v", container.Name, device)
			return nil, fmt.Errorf("cannot find volume %q to pass into container %q", device.Name, container.Name)
		}
		// Get a symbolic link associated to a block device under pod device path
		dirPath, volName := vol.BlockVolumeMapper.GetPodDeviceMapPath()
		symlinkPath := path.Join(dirPath, volName)
		if islinkExist, checkErr := blkutil.IsSymlinkExist(symlinkPath); checkErr != nil {
			return nil, checkErr
		} else if islinkExist {
			// Check readOnly in PVCVolumeSource and set read only permission if it's true.
			permission := "mrw"
			if vol.ReadOnly {
				permission = "r"
			}
			klog.V(4).Infof("Device will be attached to container %q. Path on host: %v", container.Name, symlinkPath)
			devices = append(devices, kubecontainer.DeviceInfo{PathOnHost: symlinkPath, PathInContainer: device.DevicePath, Permissions: permission})
		}
	}

	return devices, nil
}

// GeneratePodHostNameAndDomain creates a hostname and domain name for a pod,
// given that pod's spec and annotations or returns an error.
func (kl *Kubelet) GeneratePodHostNameAndDomain(pod *v1.Pod) (string, string, error) {
	clusterDomain := kl.dnsConfigurer.ClusterDomain

	hostname := pod.Name
	if len(pod.Spec.Hostname) > 0 {
		if msgs := utilvalidation.IsDNS1123Label(pod.Spec.Hostname); len(msgs) != 0 {
			return "", "", fmt.Errorf("pod Hostname %q is not a valid DNS label: %s", pod.Spec.Hostname, strings.Join(msgs, ";"))
		}
		hostname = pod.Spec.Hostname
	}

	hostname, err := truncatePodHostnameIfNeeded(pod.Name, hostname)
	if err != nil {
		return "", "", err
	}

	hostDomain := ""
	if len(pod.Spec.Subdomain) > 0 {
		if msgs := utilvalidation.IsDNS1123Label(pod.Spec.Subdomain); len(msgs) != 0 {
			return "", "", fmt.Errorf("pod Subdomain %q is not a valid DNS label: %s", pod.Spec.Subdomain, strings.Join(msgs, ";"))
		}
		hostDomain = fmt.Sprintf("%s.%s.svc.%s", pod.Spec.Subdomain, pod.Namespace, clusterDomain)
	}

	return hostname, hostDomain, nil
}

// GetPodCgroupParent gets pod cgroup parent from container manager.
func (kl *Kubelet) GetPodCgroupParent(pod *v1.Pod) string {
	pcm := kl.containerManager.NewPodContainerManager()
	_, cgroupParent := pcm.GetPodContainerName(pod)
	return cgroupParent
}

// GenerateRunContainerOptions 返回一个 opt 对象, 包含该 Pod 的 env 环境变量, 
// volume 卷列表(同时创建 hostPath 类型的目录), 端口映射信息, hostname 等.
//
// caller: pkg/kubelet/kuberuntime/kuberuntime_container.go ->
//             kubeGenericRuntimeManager.generateContainerConfig()
//
// GenerateRunContainerOptions generates the RunContainerOptions, which can be used by
// the container runtime to set parameters for launching a container.
func (kl *Kubelet) GenerateRunContainerOptions(
	pod *v1.Pod, container *v1.Container, podIP string, podIPs []string,
) (*kubecontainer.RunContainerOptions, func(), error) {
	opts, err := kl.containerManager.GetResources(pod, container)
	if err != nil {
		return nil, nil, err
	}

	hostname, hostDomainName, err := kl.GeneratePodHostNameAndDomain(pod)
	if err != nil {
		return nil, nil, err
	}
	opts.Hostname = hostname
	podName := volumeutil.GetUniquePodName(pod)
	// 从 actualStateOfWorld 中获取该 Pod 已经声明的 volume 列表(map).
	volumes := kl.volumeManager.GetMountedVolumesForPod(podName)
	// 获取端口映射配置记录(不过并没有做什么操作(比如iptables))
	opts.PortMappings = kubecontainer.MakePortMappings(container)

	// TODO: remove feature gate check after no longer needed
	if utilfeature.DefaultFeatureGate.Enabled(features.BlockVolume) {
		blkutil := volumepathhandler.NewBlockVolumePathHandler()
		blkVolumes, err := kl.makeBlockVolumes(pod, container, volumes, blkutil)
		if err != nil {
			return nil, nil, err
		}
		opts.Devices = append(opts.Devices, blkVolumes...)
	}
	// 解析环境变量
	envs, err := kl.makeEnvironmentVariables(pod, container, podIP, podIPs)
	if err != nil {
		return nil, nil, err
	}
	opts.Envs = append(opts.Envs, envs...)

	// only podIPs is sent to makeMounts, as podIPs is populated
	// even if dual-stack feature flag is not enabled.
	mounts, cleanupAction, err := makeMounts(
		pod, kl.getPodDir(pod.UID), container, hostname, hostDomainName, podIPs, 
		volumes, kl.hostutil, kl.subpather, opts.Envs,
	)
	if err != nil {
		return nil, cleanupAction, err
	}
	opts.Mounts = append(opts.Mounts, mounts...)

	// adding TerminationMessagePath on Windows is only allowed if ContainerD is used.
	// Individual files cannot be mounted as volumes using Docker for Windows.
	supportsSingleFileMapping := kl.containerRuntime.SupportsSingleFileMapping()
	if len(container.TerminationMessagePath) != 0 && supportsSingleFileMapping {
		p := kl.getPodContainerDir(pod.UID, container.Name)
		if err := os.MkdirAll(p, 0750); err != nil {
			klog.Errorf("Error on creating %q: %v", p, err)
		} else {
			opts.PodContainerDir = p
		}
	}

	// only do this check if the experimental behavior is enabled, otherwise allow it to default to false
	if kl.experimentalHostUserNamespaceDefaulting {
		opts.EnableHostUserNamespace = kl.enableHostUserNamespace(pod)
	}

	return opts, cleanupAction, nil
}

var masterServices = sets.NewString("kubernetes")

// getServiceEnvVarMap makes a map[string]string of env vars for services a
// pod in namespace ns should see.
func (kl *Kubelet) getServiceEnvVarMap(ns string, enableServiceLinks bool) (map[string]string, error) {
	var (
		serviceMap = make(map[string]*v1.Service)
		m          = make(map[string]string)
	)

	// Get all service resources from the master (via a cache),
	// and populate them into service environment variables.
	if kl.serviceLister == nil {
		// Kubelets without masters (e.g. plain GCE ContainerVM) don't set env vars.
		return m, nil
	}
	services, err := kl.serviceLister.List(labels.Everything())
	if err != nil {
		return m, fmt.Errorf("failed to list services when setting up env vars")
	}

	// project the services in namespace ns onto the master services
	for i := range services {
		service := services[i]
		// ignore services where ClusterIP is "None" or empty
		if !v1helper.IsServiceIPSet(service) {
			continue
		}
		serviceName := service.Name

		// We always want to add environment variabled for master services
		// from the master service namespace, even if enableServiceLinks is false.
		// We also add environment variables for other services in the same
		// namespace, if enableServiceLinks is true.
		if service.Namespace == kl.masterServiceNamespace && masterServices.Has(serviceName) {
			if _, exists := serviceMap[serviceName]; !exists {
				serviceMap[serviceName] = service
			}
		} else if service.Namespace == ns && enableServiceLinks {
			serviceMap[serviceName] = service
		}
	}

	mappedServices := []*v1.Service{}
	for key := range serviceMap {
		mappedServices = append(mappedServices, serviceMap[key])
	}

	for _, e := range envvars.FromServices(mappedServices) {
		m[e.Name] = e.Value
	}
	return m, nil
}

// makeEnvironmentVariables 读取Pod中的环境变量, 同时合并 configmap, secret 中解析出的环境变量.
//
// 在此函数中, 有通过 kube client 获取 configmap, secret 的操作.
//
// 注意: 此时已有 PodIP 传入, 说明 pause 已经启动, cni ipam 过程已经完成.
//
// caller: Kubelet.GenerateRunContainerOptions()
//
// Make the environment variables for a pod in the given namespace.
func (kl *Kubelet) makeEnvironmentVariables(
	pod *v1.Pod, container *v1.Container, podIP string, podIPs []string,
) ([]kubecontainer.EnvVar, error) {
	if pod.Spec.EnableServiceLinks == nil {
		return nil, fmt.Errorf(
			"nil pod.spec.enableServiceLinks encountered, cannot construct envvars",
		)
	}

	var result []kubecontainer.EnvVar
	// Note:  These are added to the docker Config, but are not included in the checksum computed
	// by kubecontainer.HashContainer(...).  That way, we can still determine whether an
	// v1.Container is already running by its hash. (We don't want to restart a container just
	// because some service changed.)
	//
	// Note that there is a race between Kubelet seeing the pod and kubelet seeing the service.
	// To avoid this users can: (1) wait between starting a service and starting; or (2) detect
	// missing service env var and exit and be restarted; or (3) use DNS instead of env vars
	// and keep trying to resolve the DNS name of the service (recommended).
	serviceEnv, err := kl.getServiceEnvVarMap(pod.Namespace, *pod.Spec.EnableServiceLinks)
	if err != nil {
		return result, err
	}

	var (
		configMaps = make(map[string]*v1.ConfigMap)
		secrets    = make(map[string]*v1.Secret)
		tmpEnv     = make(map[string]string)
	)

	// Env will override EnvFrom variables.
	// Process EnvFrom first then allow Env to replace existing values.
	for _, envFrom := range container.EnvFrom {
		switch {
		case envFrom.ConfigMapRef != nil:
			cm := envFrom.ConfigMapRef
			name := cm.Name
			configMap, ok := configMaps[name]
			if !ok {
				if kl.kubeClient == nil {
					return result, fmt.Errorf("couldn't get configMap %v/%v, no kubeClient defined", pod.Namespace, name)
				}
				optional := cm.Optional != nil && *cm.Optional
				configMap, err = kl.configMapManager.GetConfigMap(pod.Namespace, name)
				if err != nil {
					if errors.IsNotFound(err) && optional {
						// ignore error when marked optional
						continue
					}
					return result, err
				}
				configMaps[name] = configMap
			}

			invalidKeys := []string{}
			for k, v := range configMap.Data {
				if len(envFrom.Prefix) > 0 {
					k = envFrom.Prefix + k
				}
				if errMsgs := utilvalidation.IsEnvVarName(k); len(errMsgs) != 0 {
					invalidKeys = append(invalidKeys, k)
					continue
				}
				tmpEnv[k] = v
			}
			if len(invalidKeys) > 0 {
				sort.Strings(invalidKeys)
				kl.recorder.Eventf(pod, v1.EventTypeWarning, "InvalidEnvironmentVariableNames", "Keys [%s] from the EnvFrom configMap %s/%s were skipped since they are considered invalid environment variable names.", strings.Join(invalidKeys, ", "), pod.Namespace, name)
			}
		case envFrom.SecretRef != nil:
			s := envFrom.SecretRef
			name := s.Name
			secret, ok := secrets[name]
			if !ok {
				if kl.kubeClient == nil {
					return result, fmt.Errorf("couldn't get secret %v/%v, no kubeClient defined", pod.Namespace, name)
				}
				optional := s.Optional != nil && *s.Optional
				secret, err = kl.secretManager.GetSecret(pod.Namespace, name)
				if err != nil {
					if errors.IsNotFound(err) && optional {
						// ignore error when marked optional
						continue
					}
					return result, err
				}
				secrets[name] = secret
			}

			invalidKeys := []string{}
			for k, v := range secret.Data {
				if len(envFrom.Prefix) > 0 {
					k = envFrom.Prefix + k
				}
				if errMsgs := utilvalidation.IsEnvVarName(k); len(errMsgs) != 0 {
					invalidKeys = append(invalidKeys, k)
					continue
				}
				tmpEnv[k] = string(v)
			}
			if len(invalidKeys) > 0 {
				sort.Strings(invalidKeys)
				kl.recorder.Eventf(pod, v1.EventTypeWarning, "InvalidEnvironmentVariableNames", "Keys [%s] from the EnvFrom secret %s/%s were skipped since they are considered invalid environment variable names.", strings.Join(invalidKeys, ", "), pod.Namespace, name)
			}
		}
	}

	// Determine the final values of variables:
	//
	// 1.  Determine the final value of each variable:
	//     a.  If the variable's Value is set, expand the `$(var)` references to other
	//         variables in the .Value field; the sources of variables are the declared
	//         variables of the container and the service environment variables
	//     b.  If a source is defined for an environment variable, resolve the source
	// 2.  Create the container's environment in the order variables are declared
	// 3.  Add remaining service environment vars
	var (
		mappingFunc = expansion.MappingFuncFor(tmpEnv, serviceEnv)
	)
	// 注意: 上面是 envFrom, 这里是 valueFrom, 想一想 podIP 的注入方式.
	for _, envVar := range container.Env {
		runtimeVal := envVar.Value
		if runtimeVal != "" {
			// Step 1a: expand variable references
			runtimeVal = expansion.Expand(runtimeVal, mappingFunc)
		} else if envVar.ValueFrom != nil {
			// Step 1b: resolve alternate env var sources
			switch {
			case envVar.ValueFrom.FieldRef != nil:
				// 这里是从 Pod status 中注入的过程.
				runtimeVal, err = kl.podFieldSelectorRuntimeValue(
					envVar.ValueFrom.FieldRef, pod, podIP, podIPs,
				)
				if err != nil {
					return result, err
				}
			case envVar.ValueFrom.ResourceFieldRef != nil:
				defaultedPod, defaultedContainer, err := kl.defaultPodLimitsForDownwardAPI(pod, container)
				if err != nil {
					return result, err
				}
				runtimeVal, err = containerResourceRuntimeValue(
					envVar.ValueFrom.ResourceFieldRef, defaultedPod, defaultedContainer,
				)
				if err != nil {
					return result, err
				}
			case envVar.ValueFrom.ConfigMapKeyRef != nil:
				cm := envVar.ValueFrom.ConfigMapKeyRef
				name := cm.Name
				key := cm.Key
				optional := cm.Optional != nil && *cm.Optional
				configMap, ok := configMaps[name]
				if !ok {
					if kl.kubeClient == nil {
						return result, fmt.Errorf("couldn't get configMap %v/%v, no kubeClient defined", pod.Namespace, name)
					}
					configMap, err = kl.configMapManager.GetConfigMap(pod.Namespace, name)
					if err != nil {
						if errors.IsNotFound(err) && optional {
							// ignore error when marked optional
							continue
						}
						return result, err
					}
					configMaps[name] = configMap
				}
				runtimeVal, ok = configMap.Data[key]
				if !ok {
					if optional {
						continue
					}
					return result, fmt.Errorf("couldn't find key %v in ConfigMap %v/%v", key, pod.Namespace, name)
				}
			case envVar.ValueFrom.SecretKeyRef != nil:
				s := envVar.ValueFrom.SecretKeyRef
				name := s.Name
				key := s.Key
				optional := s.Optional != nil && *s.Optional
				secret, ok := secrets[name]
				if !ok {
					if kl.kubeClient == nil {
						return result, fmt.Errorf("couldn't get secret %v/%v, no kubeClient defined", pod.Namespace, name)
					}
					secret, err = kl.secretManager.GetSecret(pod.Namespace, name)
					if err != nil {
						if errors.IsNotFound(err) && optional {
							// ignore error when marked optional
							continue
						}
						return result, err
					}
					secrets[name] = secret
				}
				runtimeValBytes, ok := secret.Data[key]
				if !ok {
					if optional {
						continue
					}
					return result, fmt.Errorf("couldn't find key %v in Secret %v/%v", key, pod.Namespace, name)
				}
				runtimeVal = string(runtimeValBytes)
			}
		}
		// Accesses apiserver+Pods.
		// So, the master may set service env vars, or kubelet may.  In case both are doing
		// it, we delete the key from the kubelet-generated ones so we don't have duplicate
		// env vars.
		// TODO: remove this next line once all platforms use apiserver+Pods.
		delete(serviceEnv, envVar.Name)

		tmpEnv[envVar.Name] = runtimeVal
	}

	// Append the env vars
	for k, v := range tmpEnv {
		result = append(result, kubecontainer.EnvVar{Name: k, Value: v})
	}

	// Append remaining service env vars.
	for k, v := range serviceEnv {
		// Accesses apiserver+Pods.
		// So, the master may set service env vars, or kubelet may.  In case both are doing
		// it, we skip the key from the kubelet-generated ones so we don't have duplicate
		// env vars.
		// TODO: remove this next line once all platforms use apiserver+Pods.
		if _, present := tmpEnv[k]; !present {
			result = append(result, kubecontainer.EnvVar{Name: k, Value: v})
		}
	}
	return result, nil
}

// podFieldSelectorRuntimeValue 从 开始执行的 Pod spec, status 等部分获取可引用的环境变量.
// 不过要注意, 这种方式可引用的字段是有限的.
//
// caller: Kubelet.makeEnvironmentVariables()
//
// podFieldSelectorRuntimeValue returns the runtime value of the given
// selector for a pod.
func (kl *Kubelet) podFieldSelectorRuntimeValue(
	fs *v1.ObjectFieldSelector, pod *v1.Pod, podIP string, podIPs []string,
) (string, error) {
	internalFieldPath, _, err := podshelper.ConvertDownwardAPIFieldLabel(fs.APIVersion, fs.FieldPath, "")
	if err != nil {
		return "", err
	}
	switch internalFieldPath {
	case "spec.nodeName":
		return pod.Spec.NodeName, nil
	case "spec.serviceAccountName":
		return pod.Spec.ServiceAccountName, nil
	case "status.hostIP":
		hostIP, err := kl.getHostIPAnyWay()
		if err != nil {
			return "", err
		}
		return hostIP.String(), nil
	case "status.podIP":
		return podIP, nil
	case "status.podIPs":
		return strings.Join(podIPs, ","), nil
	}
	return fieldpath.ExtractFieldPathAsString(pod, internalFieldPath)
}

// containerResourceRuntimeValue returns the value of the provided container resource
func containerResourceRuntimeValue(fs *v1.ResourceFieldSelector, pod *v1.Pod, container *v1.Container) (string, error) {
	containerName := fs.ContainerName
	if len(containerName) == 0 {
		return resource.ExtractContainerResourceValue(fs, container)
	}
	return resource.ExtractResourceValueByContainerName(fs, pod, containerName)
}

// caller:
// 	1. pkg/kubelet/kubelet.go -> Kubelet.syncPod()
//
// One of the following arguments must be non-nil: runningPod, status.
// TODO: Modify containerRuntime.KillPod() to accept the right arguments.
func (kl *Kubelet) killPod(
	pod *v1.Pod, runningPod *kubecontainer.Pod, status *kubecontainer.PodStatus, 
	gracePeriodOverride *int64,
) error {
	var p kubecontainer.Pod
	// runningPod 与 status 不可同时为 nil
	if runningPod != nil {
		p = *runningPod
	} else if status != nil {
		p = kubecontainer.ConvertPodStatusToRunningPod(kl.getRuntime().Type(), status)
	} else {
		return fmt.Errorf("one of the two arguments must be non-nil: runningPod, status")
	}

	// Call the container runtime KillPod method which stops all running containers of the pod
	if err := kl.containerRuntime.KillPod(pod, p, gracePeriodOverride); err != nil {
		return err
	}
	if err := kl.containerManager.UpdateQOSCgroups(); err != nil {
		klog.V(2).Infof("Failed to update QoS cgroups while killing pod: %v", err)
	}
	return nil
}

// makePodDataDirs 分别在 /var/lib/kubelet 的 pods, volumes, plugins 目录下,
// 创建名为 pod uid 的子目录.
//
// caller: pkg/kubelet/kubelet.go -> Kubelet.syncPod()
//
// makePodDataDirs creates the dirs for the pod datas.
func (kl *Kubelet) makePodDataDirs(pod *v1.Pod) error {
	uid := pod.UID
	if err := os.MkdirAll(kl.getPodDir(uid), 0750); err != nil && !os.IsExist(err) {
		return err
	}
	if err := os.MkdirAll(kl.getPodVolumesDir(uid), 0750); err != nil && !os.IsExist(err) {
		return err
	}
	if err := os.MkdirAll(kl.getPodPluginsDir(uid), 0750); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// getPullSecretsForPod inspects the Pod and retrieves the referenced pull
// secrets.
func (kl *Kubelet) getPullSecretsForPod(pod *v1.Pod) []v1.Secret {
	pullSecrets := []v1.Secret{}

	for _, secretRef := range pod.Spec.ImagePullSecrets {
		secret, err := kl.secretManager.GetSecret(pod.Namespace, secretRef.Name)
		if err != nil {
			klog.Warningf("Unable to retrieve pull secret %s/%s for %s/%s due to %v.  The image pull may not succeed.", pod.Namespace, secretRef.Name, pod.Namespace, pod.Name, err)
			continue
		}

		pullSecrets = append(pullSecrets, *secret)
	}

	return pullSecrets
}

// podIsTerminated returns true if pod is in the terminated state ("Failed" or "Succeeded").
func (kl *Kubelet) podIsTerminated(pod *v1.Pod) bool {
	// Check the cached pod status which was set after the last sync.
	status, ok := kl.statusManager.GetPodStatus(pod.UID)
	if !ok {
		// If there is no cached status, use the status from the
		// apiserver. This is useful if kubelet has recently been
		// restarted.
		status = pod.Status
	}
	return status.Phase == v1.PodFailed || status.Phase == v1.PodSucceeded || (pod.DeletionTimestamp != nil && notRunning(status.ContainerStatuses))
}

// IsPodTerminated returns true if the pod with the provided UID is in a terminated state ("Failed" or "Succeeded")
// or if the pod has been deleted or removed
func (kl *Kubelet) IsPodTerminated(uid types.UID) bool {
	pod, podFound := kl.podManager.GetPodByUID(uid)
	if !podFound {
		return true
	}
	return kl.podIsTerminated(pod)
}

// IsPodDeleted returns true if the pod is deleted.  For the pod to be deleted, either:
// 1. The pod object is deleted
// 2. The pod's status is evicted
// 3. The pod's deletion timestamp is set, and containers are not running
func (kl *Kubelet) IsPodDeleted(uid types.UID) bool {
	pod, podFound := kl.podManager.GetPodByUID(uid)
	if !podFound {
		return true
	}
	status, statusFound := kl.statusManager.GetPodStatus(pod.UID)
	if !statusFound {
		status = pod.Status
	}
	return eviction.PodIsEvicted(status) || (pod.DeletionTimestamp != nil && notRunning(status.ContainerStatuses))
}

// PodResourcesAreReclaimed returns true if all required node-level resources
// that a pod was consuming have been reclaimed by the kubelet. 
// Reclaiming resources is a prerequisite to deleting a pod from the API server.
func (kl *Kubelet) PodResourcesAreReclaimed(pod *v1.Pod, status v1.PodStatus) bool {
	if !notRunning(status.ContainerStatuses) {
		// We shouldn't delete pods that still have running containers
		klog.V(3).Infof("Pod %q is terminated, but some containers are still running", format.Pod(pod))
		return false
	}
	// pod's containers should be deleted
	runtimeStatus, err := kl.podCache.Get(pod.UID)
	if err != nil {
		klog.V(3).Infof("Pod %q is terminated, Error getting runtimeStatus from the podCache: %s", format.Pod(pod), err)
		return false
	}
	if len(runtimeStatus.ContainerStatuses) > 0 {
		var statusStr string
		for _, status := range runtimeStatus.ContainerStatuses {
			statusStr += fmt.Sprintf("%+v ", *status)
		}
		klog.V(3).Infof("Pod %q is terminated, but some containers have not been cleaned up: %s", format.Pod(pod), statusStr)
		return false
	}
	if kl.podVolumesExist(pod.UID) && !kl.keepTerminatedPodVolumes {
		// We shouldn't delete pods whose volumes have not been cleaned up if we are not keeping terminated pod volumes
		klog.V(3).Infof("Pod %q is terminated, but some volumes have not been cleaned up", format.Pod(pod))
		return false
	}
	if kl.kubeletConfiguration.CgroupsPerQOS {
		pcm := kl.containerManager.NewPodContainerManager()
		if pcm.Exists(pod) {
			klog.V(3).Infof("Pod %q is terminated, but pod cgroup sandbox has not been cleaned up", format.Pod(pod))
			return false
		}
	}
	return true
}

// podResourcesAreReclaimed simply calls PodResourcesAreReclaimed with the most up-to-date status.
func (kl *Kubelet) podResourcesAreReclaimed(pod *v1.Pod) bool {
	status, ok := kl.statusManager.GetPodStatus(pod.UID)
	if !ok {
		status = pod.Status
	}
	return kl.PodResourcesAreReclaimed(pod, status)
}

// notRunning returns true if every status is terminated or waiting, or the status list
// is empty.
func notRunning(statuses []v1.ContainerStatus) bool {
	for _, status := range statuses {
		if status.State.Terminated == nil && status.State.Waiting == nil {
			return false
		}
	}
	return true
}

// filterOutTerminatedPods returns the given pods which the status manager
// does not consider failed or succeeded.
func (kl *Kubelet) filterOutTerminatedPods(pods []*v1.Pod) []*v1.Pod {
	var filteredPods []*v1.Pod
	for _, p := range pods {
		if kl.podIsTerminated(p) {
			continue
		}
		filteredPods = append(filteredPods, p)
	}
	return filteredPods
}

// removeOrphanedPodStatuses removes obsolete entries in podStatus where
// the pod is no longer considered bound to this node.
func (kl *Kubelet) removeOrphanedPodStatuses(pods []*v1.Pod, mirrorPods []*v1.Pod) {
	podUIDs := make(map[types.UID]bool)
	for _, pod := range pods {
		podUIDs[pod.UID] = true
	}
	for _, pod := range mirrorPods {
		podUIDs[pod.UID] = true
	}
	kl.statusManager.RemoveOrphanedStatuses(podUIDs)
}

// HandlePodCleanups performs a series of cleanup work, including terminating
// pod workers, killing unwanted pods, and removing orphaned volumes/pod
// directories.
// NOTE: This function is executed by the main sync loop, so it
// should not contain any blocking calls.
func (kl *Kubelet) HandlePodCleanups() error {
	// The kubelet lacks checkpointing, so we need to introspect the set of pods
	// in the cgroup tree prior to inspecting the set of pods in our pod manager.
	// this ensures our view of the cgroup tree does not mistakenly observe pods
	// that are added after the fact...
	var (
		cgroupPods map[types.UID]cm.CgroupName
		err        error
	)
	if kl.cgroupsPerQOS {
		pcm := kl.containerManager.NewPodContainerManager()
		cgroupPods, err = pcm.GetAllPodsFromCgroups()
		if err != nil {
			return fmt.Errorf("failed to get list of pods that still exist on cgroup mounts: %v", err)
		}
	}

	allPods, mirrorPods := kl.podManager.GetPodsAndMirrorPods()
	// Pod phase progresses monotonically. Once a pod has reached a final state,
	// it should never leave regardless of the restart policy. The statuses
	// of such pods should not be changed, and there is no need to sync them.
	// TODO: the logic here does not handle two cases:
	//   1. If the containers were removed immediately after they died, kubelet
	//      may fail to generate correct statuses, let alone filtering correctly.
	//   2. If kubelet restarted before writing the terminated status for a pod
	//      to the apiserver, it could still restart the terminated pod (even
	//      though the pod was not considered terminated by the apiserver).
	// These two conditions could be alleviated by checkpointing kubelet.
	activePods := kl.filterOutTerminatedPods(allPods)

	desiredPods := make(map[types.UID]sets.Empty)
	for _, pod := range activePods {
		desiredPods[pod.UID] = sets.Empty{}
	}
	// Stop the workers for no-longer existing pods.
	// TODO: is here the best place to forget pod workers?
	kl.podWorkers.ForgetNonExistingPodWorkers(desiredPods)
	kl.probeManager.CleanupPods(desiredPods)

	runningPods, err := kl.runtimeCache.GetPods()
	if err != nil {
		klog.Errorf("Error listing containers: %#v", err)
		return err
	}
	for _, pod := range runningPods {
		if _, found := desiredPods[pod.ID]; !found {
			kl.podKillingCh <- &kubecontainer.PodPair{APIPod: nil, RunningPod: pod}
		}
	}

	kl.removeOrphanedPodStatuses(allPods, mirrorPods)
	// Note that we just killed the unwanted pods. This may not have reflected
	// in the cache. We need to bypass the cache to get the latest set of
	// running pods to clean up the volumes.
	// TODO: Evaluate the performance impact of bypassing the runtime cache.
	runningPods, err = kl.containerRuntime.GetPods(false)
	if err != nil {
		klog.Errorf("Error listing containers: %#v", err)
		return err
	}

	// Remove any orphaned volumes.
	// Note that we pass all pods (including terminated pods) to the function,
	// so that we don't remove volumes associated with terminated but not yet
	// deleted pods.
	err = kl.cleanupOrphanedPodDirs(allPods, runningPods)
	if err != nil {
		// We want all cleanup tasks to be run even if one of them failed. So
		// we just log an error here and continue other cleanup tasks.
		// This also applies to the other clean up tasks.
		klog.Errorf("Failed cleaning up orphaned pod directories: %v", err)
	}

	// Remove any orphaned mirror pods.
	kl.podManager.DeleteOrphanedMirrorPods()

	// Remove any cgroups in the hierarchy for pods that are no longer running.
	if kl.cgroupsPerQOS {
		kl.cleanupOrphanedPodCgroups(cgroupPods, activePods)
	}

	kl.backOff.GC()
	return nil
}

// Kubelet.podKillingCh 成员在 NewMainKubelet() 函数中被初始化,
// 在 Kubelet.deletePod() 函数中写入数据.
//
// podKiller launches a goroutine to kill a pod received from the channel if
// another goroutine isn't already in action.
func (kl *Kubelet) podKiller() {
	killing := sets.NewString()
	// guard for the killing set
	lock := sync.Mutex{}
	for podPair := range kl.podKillingCh {
		runningPod := podPair.RunningPod
		apiPod := podPair.APIPod

		lock.Lock()
		exists := killing.Has(string(runningPod.ID))
		if !exists {
			killing.Insert(string(runningPod.ID))
		}
		lock.Unlock()

		if !exists {
			go func(apiPod *v1.Pod, runningPod *kubecontainer.Pod) {
				klog.V(2).Infof("Killing unwanted pod %q", runningPod.Name)
				err := kl.killPod(apiPod, runningPod, nil, nil)
				if err != nil {
					klog.Errorf("Failed killing the pod %q: %v", runningPod.Name, err)
				}
				lock.Lock()
				killing.Delete(string(runningPod.ID))
				lock.Unlock()
			}(apiPod, runningPod)
		}
	}
}

// validateContainerLogStatus returns the container ID for the desired container to retrieve logs for, based on the state
// of the container. The previous flag will only return the logs for the last terminated container, otherwise, the current
// running container is preferred over a previous termination. If info about the container is not available then a specific
// error is returned to the end user.
func (kl *Kubelet) validateContainerLogStatus(podName string, podStatus *v1.PodStatus, containerName string, previous bool) (containerID kubecontainer.ContainerID, err error) {
	var cID string

	cStatus, found := podutil.GetContainerStatus(podStatus.ContainerStatuses, containerName)
	if !found {
		cStatus, found = podutil.GetContainerStatus(podStatus.InitContainerStatuses, containerName)
	}
	if !found && utilfeature.DefaultFeatureGate.Enabled(features.EphemeralContainers) {
		cStatus, found = podutil.GetContainerStatus(podStatus.EphemeralContainerStatuses, containerName)
	}
	if !found {
		return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is not available", containerName, podName)
	}
	lastState := cStatus.LastTerminationState
	waiting, running, terminated := cStatus.State.Waiting, cStatus.State.Running, cStatus.State.Terminated

	switch {
	case previous:
		if lastState.Terminated == nil || lastState.Terminated.ContainerID == "" {
			return kubecontainer.ContainerID{}, fmt.Errorf("previous terminated container %q in pod %q not found", containerName, podName)
		}
		cID = lastState.Terminated.ContainerID

	case running != nil:
		cID = cStatus.ContainerID

	case terminated != nil:
		// in cases where the next container didn't start, terminated.ContainerID will be empty, so get logs from the lastState.Terminated.
		if terminated.ContainerID == "" {
			if lastState.Terminated != nil && lastState.Terminated.ContainerID != "" {
				cID = lastState.Terminated.ContainerID
			} else {
				return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is terminated", containerName, podName)
			}
		} else {
			cID = terminated.ContainerID
		}

	case lastState.Terminated != nil:
		if lastState.Terminated.ContainerID == "" {
			return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is terminated", containerName, podName)
		}
		cID = lastState.Terminated.ContainerID

	case waiting != nil:
		// output some info for the most common pending failures
		switch reason := waiting.Reason; reason {
		case images.ErrImagePull.Error():
			return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is waiting to start: image can't be pulled", containerName, podName)
		case images.ErrImagePullBackOff.Error():
			return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is waiting to start: trying and failing to pull image", containerName, podName)
		default:
			return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is waiting to start: %v", containerName, podName, reason)
		}
	default:
		// unrecognized state
		return kubecontainer.ContainerID{}, fmt.Errorf("container %q in pod %q is waiting to start - no logs yet", containerName, podName)
	}

	return kubecontainer.ParseContainerID(cID), nil
}

// GetKubeletContainerLogs returns logs from the container
// TODO: this method is returning logs of random container attempts, when it should be returning the most recent attempt
// or all of them.
func (kl *Kubelet) GetKubeletContainerLogs(ctx context.Context, podFullName, containerName string, logOptions *v1.PodLogOptions, stdout, stderr io.Writer) error {
	// Pod workers periodically write status to statusManager. If status is not
	// cached there, something is wrong (or kubelet just restarted and hasn't
	// caught up yet). Just assume the pod is not ready yet.
	name, namespace, err := kubecontainer.ParsePodFullName(podFullName)
	if err != nil {
		return fmt.Errorf("unable to parse pod full name %q: %v", podFullName, err)
	}

	pod, ok := kl.GetPodByName(namespace, name)
	if !ok {
		return fmt.Errorf("pod %q cannot be found - no logs available", name)
	}

	podUID := pod.UID
	if mirrorPod, ok := kl.podManager.GetMirrorPodByPod(pod); ok {
		podUID = mirrorPod.UID
	}
	podStatus, found := kl.statusManager.GetPodStatus(podUID)
	if !found {
		// If there is no cached status, use the status from the
		// apiserver. This is useful if kubelet has recently been
		// restarted.
		podStatus = pod.Status
	}

	// TODO: Consolidate the logic here with kuberuntime.GetContainerLogs, here we convert container name to containerID,
	// but inside kuberuntime we convert container id back to container name and restart count.
	// TODO: After separate container log lifecycle management, we should get log based on the existing log files
	// instead of container status.
	containerID, err := kl.validateContainerLogStatus(pod.Name, &podStatus, containerName, logOptions.Previous)
	if err != nil {
		return err
	}

	// Do a zero-byte write to stdout before handing off to the container runtime.
	// This ensures at least one Write call is made to the writer when copying starts,
	// even if we then block waiting for log output from the container.
	if _, err := stdout.Write([]byte{}); err != nil {
		return err
	}

	if kl.dockerLegacyService != nil {
		// dockerLegacyService should only be non-nil when we actually need it, so
		// inject it into the runtimeService.
		// TODO(random-liu): Remove this hack after deprecating unsupported log driver.
		return kl.dockerLegacyService.GetContainerLogs(ctx, pod, containerID, logOptions, stdout, stderr)
	}
	return kl.containerRuntime.GetContainerLogs(ctx, pod, containerID, logOptions, stdout, stderr)
}

// caller:
// 	1. Kubelet.generateAPIPodStatus() 只有这一处
//
// getPhase returns the phase of a pod given its container info.
func getPhase(spec *v1.PodSpec, info []v1.ContainerStatus) v1.PodPhase {
	pendingInitialization := 0
	failedInitialization := 0
	for _, container := range spec.InitContainers {
		containerStatus, ok := podutil.GetContainerStatus(info, container.Name)
		if !ok {
			pendingInitialization++
			continue
		}

		switch {
		case containerStatus.State.Running != nil:
			pendingInitialization++
		case containerStatus.State.Terminated != nil:
			if containerStatus.State.Terminated.ExitCode != 0 {
				failedInitialization++
			}
		case containerStatus.State.Waiting != nil:
			if containerStatus.LastTerminationState.Terminated != nil {
				if containerStatus.LastTerminationState.Terminated.ExitCode != 0 {
					failedInitialization++
				}
			} else {
				pendingInitialization++
			}
		default:
			pendingInitialization++
		}
	}

	unknown := 0
	running := 0
	waiting := 0
	stopped := 0
	succeeded := 0
	for _, container := range spec.Containers {
		containerStatus, ok := podutil.GetContainerStatus(info, container.Name)
		if !ok {
			unknown++
			continue
		}

		switch {
		case containerStatus.State.Running != nil:
			running++
		case containerStatus.State.Terminated != nil:
			stopped++
			if containerStatus.State.Terminated.ExitCode == 0 {
				succeeded++
			}
		case containerStatus.State.Waiting != nil:
			if containerStatus.LastTerminationState.Terminated != nil {
				stopped++
			} else {
				waiting++
			}
		default:
			unknown++
		}
	}

	if failedInitialization > 0 && spec.RestartPolicy == v1.RestartPolicyNever {
		return v1.PodFailed
	}

	switch {
	case pendingInitialization > 0:
		fallthrough
	case waiting > 0:
		klog.V(5).Infof("pod waiting > 0, pending")
		// One or more containers has not been started
		return v1.PodPending
	case running > 0 && unknown == 0:
		// All containers have been started, and at least
		// one container is running
		return v1.PodRunning
	case running == 0 && stopped > 0 && unknown == 0:
		// All containers are terminated
		if spec.RestartPolicy == v1.RestartPolicyAlways {
			// All containers are in the process of restarting
			return v1.PodRunning
		}
		if stopped == succeeded {
			// RestartPolicy is not Always, and all
			// containers are terminated in success
			return v1.PodSucceeded
		}
		if spec.RestartPolicy == v1.RestartPolicyNever {
			// RestartPolicy is Never, and all containers are
			// terminated with at least one in failure
			return v1.PodFailed
		}
		// RestartPolicy is OnFailure, and at least one in failure
		// and in the process of restarting
		return v1.PodRunning
	default:
		klog.V(5).Infof("pod default case, pending")
		return v1.PodPending
	}
}

// generateAPIPodStatus creates the final API pod status for a pod, given the
// internal pod status.
func (kl *Kubelet) generateAPIPodStatus(
	pod *v1.Pod, podStatus *kubecontainer.PodStatus,
) v1.PodStatus {
	klog.V(3).Infof("Generating status for %q", format.Pod(pod))

	s := kl.convertStatusToAPIStatus(pod, podStatus)

	// check if an internal module has requested the pod is evicted.
	for _, podSyncHandler := range kl.PodSyncHandlers {
		if result := podSyncHandler.ShouldEvict(pod); result.Evict {
			s.Phase = v1.PodFailed
			s.Reason = result.Reason
			s.Message = result.Message
			return *s
		}
	}

	// Assume info is ready to process
	spec := &pod.Spec
	allStatus := append(append([]v1.ContainerStatus{}, s.ContainerStatuses...), s.InitContainerStatuses...)
	s.Phase = getPhase(spec, allStatus)
	// Check for illegal phase transition
	if pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodSucceeded {
		// API server shows terminal phase; transitions are not allowed
		if s.Phase != pod.Status.Phase {
			klog.Errorf("Pod attempted illegal phase transition from %s to %s: %v", pod.Status.Phase, s.Phase, s)
			// Force back to phase from the API server
			s.Phase = pod.Status.Phase
		}
	}
	kl.probeManager.UpdatePodStatus(pod.UID, s)
	s.Conditions = append(s.Conditions, status.GeneratePodInitializedCondition(spec, s.InitContainerStatuses, s.Phase))
	s.Conditions = append(s.Conditions, status.GeneratePodReadyCondition(spec, s.Conditions, s.ContainerStatuses, s.Phase))
	s.Conditions = append(s.Conditions, status.GenerateContainersReadyCondition(spec, s.ContainerStatuses, s.Phase))
	// Status manager will take care of the LastTransitionTimestamp, either preserve
	// the timestamp from apiserver, or set a new one. When kubelet sees the pod,
	// `PodScheduled` condition must be true.
	s.Conditions = append(s.Conditions, v1.PodCondition{
		Type:   v1.PodScheduled,
		Status: v1.ConditionTrue,
	})

	if kl.kubeClient != nil {
		hostIP, err := kl.getHostIPAnyWay()
		if err != nil {
			klog.V(4).Infof("Cannot get host IP: %v", err)
		} else {
			s.HostIP = hostIP.String()
			if kubecontainer.IsHostNetworkPod(pod) && s.PodIP == "" {
				s.PodIP = hostIP.String()
			}
		}
	}

	return *s
}

// convertStatusToAPIStatus creates an api PodStatus for the given pod from
// the given internal pod status.  It is purely transformative and does not
// alter the kubelet state at all.
func (kl *Kubelet) convertStatusToAPIStatus(pod *v1.Pod, podStatus *kubecontainer.PodStatus) *v1.PodStatus {
	var apiPodStatus v1.PodStatus
	apiPodStatus.PodIPs = make([]v1.PodIP, 0, len(podStatus.IPs))
	for _, ip := range podStatus.IPs {
		apiPodStatus.PodIPs = append(apiPodStatus.PodIPs, v1.PodIP{
			IP: ip,
		})
	}

	if len(apiPodStatus.PodIPs) > 0 {
		apiPodStatus.PodIP = apiPodStatus.PodIPs[0].IP
	}
	// set status for Pods created on versions of kube older than 1.6
	apiPodStatus.QOSClass = v1qos.GetPodQOS(pod)

	oldPodStatus, found := kl.statusManager.GetPodStatus(pod.UID)
	if !found {
		oldPodStatus = pod.Status
	}

	apiPodStatus.ContainerStatuses = kl.convertToAPIContainerStatuses(
		pod, podStatus,
		oldPodStatus.ContainerStatuses,
		pod.Spec.Containers,
		len(pod.Spec.InitContainers) > 0,
		false,
	)
	apiPodStatus.InitContainerStatuses = kl.convertToAPIContainerStatuses(
		pod, podStatus,
		oldPodStatus.InitContainerStatuses,
		pod.Spec.InitContainers,
		len(pod.Spec.InitContainers) > 0,
		true,
	)
	if utilfeature.DefaultFeatureGate.Enabled(features.EphemeralContainers) {
		var ecSpecs []v1.Container
		for i := range pod.Spec.EphemeralContainers {
			ecSpecs = append(ecSpecs, v1.Container(pod.Spec.EphemeralContainers[i].EphemeralContainerCommon))
		}

		// #80875: By now we've iterated podStatus 3 times. We could refactor this to make a single
		// pass through podStatus.ContainerStatuses
		apiPodStatus.EphemeralContainerStatuses = kl.convertToAPIContainerStatuses(
			pod, podStatus,
			oldPodStatus.EphemeralContainerStatuses,
			ecSpecs,
			len(pod.Spec.InitContainers) > 0,
			false,
		)
	}

	// Preserves conditions not controlled by kubelet
	for _, c := range pod.Status.Conditions {
		if !kubetypes.PodConditionByKubelet(c.Type) {
			apiPodStatus.Conditions = append(apiPodStatus.Conditions, c)
		}
	}

	return &apiPodStatus
}

// convertToAPIContainerStatuses converts the given internal container
// statuses into API container statuses.
func (kl *Kubelet) convertToAPIContainerStatuses(pod *v1.Pod, podStatus *kubecontainer.PodStatus, previousStatus []v1.ContainerStatus, containers []v1.Container, hasInitContainers, isInitContainer bool) []v1.ContainerStatus {
	convertContainerStatus := func(cs *kubecontainer.ContainerStatus) *v1.ContainerStatus {
		cid := cs.ID.String()
		status := &v1.ContainerStatus{
			Name:         cs.Name,
			RestartCount: int32(cs.RestartCount),
			Image:        cs.Image,
			ImageID:      cs.ImageID,
			ContainerID:  cid,
		}
		switch cs.State {
		case kubecontainer.ContainerStateRunning:
			status.State.Running = &v1.ContainerStateRunning{StartedAt: metav1.NewTime(cs.StartedAt)}
		case kubecontainer.ContainerStateCreated:
			// Treat containers in the "created" state as if they are exited.
			// The pod workers are supposed start all containers it creates in
			// one sync (syncPod) iteration. There should not be any normal
			// "created" containers when the pod worker generates the status at
			// the beginning of a sync iteration.
			fallthrough
		case kubecontainer.ContainerStateExited:
			status.State.Terminated = &v1.ContainerStateTerminated{
				ExitCode:    int32(cs.ExitCode),
				Reason:      cs.Reason,
				Message:     cs.Message,
				StartedAt:   metav1.NewTime(cs.StartedAt),
				FinishedAt:  metav1.NewTime(cs.FinishedAt),
				ContainerID: cid,
			}
		default:
			status.State.Waiting = &v1.ContainerStateWaiting{}
		}
		return status
	}

	// Fetch old containers statuses from old pod status.
	oldStatuses := make(map[string]v1.ContainerStatus, len(containers))
	for _, status := range previousStatus {
		oldStatuses[status.Name] = status
	}

	// Set all container statuses to default waiting state
	statuses := make(map[string]*v1.ContainerStatus, len(containers))
	defaultWaitingState := v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "ContainerCreating"}}
	if hasInitContainers {
		defaultWaitingState = v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "PodInitializing"}}
	}

	for _, container := range containers {
		status := &v1.ContainerStatus{
			Name:  container.Name,
			Image: container.Image,
			State: defaultWaitingState,
		}
		oldStatus, found := oldStatuses[container.Name]
		if found {
			if oldStatus.State.Terminated != nil {
				// Do not update status on terminated init containers as
				// they be removed at any time.
				status = &oldStatus
			} else {
				// Apply some values from the old statuses as the default values.
				status.RestartCount = oldStatus.RestartCount
				status.LastTerminationState = oldStatus.LastTerminationState
			}
		}
		statuses[container.Name] = status
	}

	// Make the latest container status comes first.
	sort.Sort(sort.Reverse(kubecontainer.SortContainerStatusesByCreationTime(podStatus.ContainerStatuses)))
	// Set container statuses according to the statuses seen in pod status
	containerSeen := map[string]int{}
	for _, cStatus := range podStatus.ContainerStatuses {
		cName := cStatus.Name
		if _, ok := statuses[cName]; !ok {
			// This would also ignore the infra container.
			continue
		}
		if containerSeen[cName] >= 2 {
			continue
		}
		status := convertContainerStatus(cStatus)
		if containerSeen[cName] == 0 {
			statuses[cName] = status
		} else {
			statuses[cName].LastTerminationState = status.State
		}
		containerSeen[cName] = containerSeen[cName] + 1
	}

	// Handle the containers failed to be started, which should be in Waiting state.
	for _, container := range containers {
		if isInitContainer {
			// If the init container is terminated with exit code 0, it won't be restarted.
			// TODO(random-liu): Handle this in a cleaner way.
			s := podStatus.FindContainerStatusByName(container.Name)
			if s != nil && s.State == kubecontainer.ContainerStateExited && s.ExitCode == 0 {
				continue
			}
		}
		// If a container should be restarted in next syncpod, it is *Waiting*.
		if !kubecontainer.ShouldContainerBeRestarted(&container, pod, podStatus) {
			continue
		}
		status := statuses[container.Name]
		reason, ok := kl.reasonCache.Get(pod.UID, container.Name)
		if !ok {
			// In fact, we could also apply Waiting state here, but it is less informative,
			// and the container will be restarted soon, so we prefer the original state here.
			// Note that with the current implementation of ShouldContainerBeRestarted the original state here
			// could be:
			//   * Waiting: There is no associated historical container and start failure reason record.
			//   * Terminated: The container is terminated.
			continue
		}
		if status.State.Terminated != nil {
			status.LastTerminationState = status.State
		}
		status.State = v1.ContainerState{
			Waiting: &v1.ContainerStateWaiting{
				Reason:  reason.Err.Error(),
				Message: reason.Message,
			},
		}
		statuses[container.Name] = status
	}

	var containerStatuses []v1.ContainerStatus
	for _, status := range statuses {
		containerStatuses = append(containerStatuses, *status)
	}

	// Sort the container statuses since clients of this interface expect the list
	// of containers in a pod has a deterministic order.
	if isInitContainer {
		kubetypes.SortInitContainerStatuses(pod, containerStatuses)
	} else {
		sort.Sort(kubetypes.SortedContainerStatuses(containerStatuses))
	}
	return containerStatuses
}

// ServeLogs returns logs of current machine.
func (kl *Kubelet) ServeLogs(w http.ResponseWriter, req *http.Request) {
	// TODO: whitelist logs we are willing to serve
	kl.logServer.ServeHTTP(w, req)
}

// findContainer finds and returns the container with the given pod ID, full name, and container name.
// It returns nil if not found.
func (kl *Kubelet) findContainer(podFullName string, podUID types.UID, containerName string) (*kubecontainer.Container, error) {
	pods, err := kl.containerRuntime.GetPods(false)
	if err != nil {
		return nil, err
	}
	// Resolve and type convert back again.
	// We need the static pod UID but the kubecontainer API works with types.UID.
	podUID = types.UID(kl.podManager.TranslatePodUID(podUID))
	pod := kubecontainer.Pods(pods).FindPod(podFullName, podUID)
	return pod.FindContainerByName(containerName), nil
}

// RunInContainer runs a command in a container, returns the combined stdout, stderr as an array of bytes
func (kl *Kubelet) RunInContainer(podFullName string, podUID types.UID, containerName string, cmd []string) ([]byte, error) {
	container, err := kl.findContainer(podFullName, podUID, containerName)
	if err != nil {
		return nil, err
	}
	if container == nil {
		return nil, fmt.Errorf("container not found (%q)", containerName)
	}
	// TODO(tallclair): Pass a proper timeout value.
	return kl.runner.RunInContainer(container.ID, cmd, 0)
}

// GetExec gets the URL the exec will be served from, or nil if the Kubelet will serve it.
func (kl *Kubelet) GetExec(podFullName string, podUID types.UID, containerName string, cmd []string, streamOpts remotecommandserver.Options) (*url.URL, error) {
	container, err := kl.findContainer(podFullName, podUID, containerName)
	if err != nil {
		return nil, err
	}
	if container == nil {
		return nil, fmt.Errorf("container not found (%q)", containerName)
	}
	return kl.streamingRuntime.GetExec(container.ID, cmd, streamOpts.Stdin, streamOpts.Stdout, streamOpts.Stderr, streamOpts.TTY)
}

// GetAttach gets the URL the attach will be served from, or nil if the Kubelet will serve it.
func (kl *Kubelet) GetAttach(podFullName string, podUID types.UID, containerName string, streamOpts remotecommandserver.Options) (*url.URL, error) {
	container, err := kl.findContainer(podFullName, podUID, containerName)
	if err != nil {
		return nil, err
	}
	if container == nil {
		return nil, fmt.Errorf("container %s not found in pod %s", containerName, podFullName)
	}

	// The TTY setting for attach must match the TTY setting in the initial container configuration,
	// since whether the process is running in a TTY cannot be changed after it has started.  We
	// need the api.Pod to get the TTY status.
	pod, found := kl.GetPodByFullName(podFullName)
	if !found || (string(podUID) != "" && pod.UID != podUID) {
		return nil, fmt.Errorf("pod %s not found", podFullName)
	}
	containerSpec := kubecontainer.GetContainerSpec(pod, containerName)
	if containerSpec == nil {
		return nil, fmt.Errorf("container %s not found in pod %s", containerName, podFullName)
	}
	tty := containerSpec.TTY

	return kl.streamingRuntime.GetAttach(container.ID, streamOpts.Stdin, streamOpts.Stdout, streamOpts.Stderr, tty)
}

// GetPortForward gets the URL the port-forward will be served from, or nil if the Kubelet will serve it.
func (kl *Kubelet) GetPortForward(podName, podNamespace string, podUID types.UID, portForwardOpts portforward.V4Options) (*url.URL, error) {
	pods, err := kl.containerRuntime.GetPods(false)
	if err != nil {
		return nil, err
	}
	// Resolve and type convert back again.
	// We need the static pod UID but the kubecontainer API works with types.UID.
	podUID = types.UID(kl.podManager.TranslatePodUID(podUID))
	podFullName := kubecontainer.BuildPodFullName(podName, podNamespace)
	pod := kubecontainer.Pods(pods).FindPod(podFullName, podUID)
	if pod.IsEmpty() {
		return nil, fmt.Errorf("pod not found (%q)", podFullName)
	}

	return kl.streamingRuntime.GetPortForward(podName, podNamespace, podUID, portForwardOpts.Ports)
}

// cleanupOrphanedPodCgroups removes cgroups that should no longer exist.
// it reconciles the cached state of cgroupPods with the specified list of runningPods
func (kl *Kubelet) cleanupOrphanedPodCgroups(cgroupPods map[types.UID]cm.CgroupName, activePods []*v1.Pod) {
	// Add all running pods to the set that we want to preserve
	podSet := sets.NewString()
	for _, pod := range activePods {
		podSet.Insert(string(pod.UID))
	}
	pcm := kl.containerManager.NewPodContainerManager()

	// Iterate over all the found pods to verify if they should be running
	for uid, val := range cgroupPods {
		// if the pod is in the running set, its not a candidate for cleanup
		if podSet.Has(string(uid)) {
			continue
		}

		// If volumes have not been unmounted/detached, do not delete the cgroup
		// so any memory backed volumes don't have their charges propagated to the
		// parent croup.  If the volumes still exist, reduce the cpu shares for any
		// process in the cgroup to the minimum value while we wait.  if the kubelet
		// is configured to keep terminated volumes, we will delete the cgroup and not block.
		if podVolumesExist := kl.podVolumesExist(uid); podVolumesExist && !kl.keepTerminatedPodVolumes {
			klog.V(3).Infof("Orphaned pod %q found, but volumes not yet removed.  Reducing cpu to minimum", uid)
			if err := pcm.ReduceCPULimits(val); err != nil {
				klog.Warningf("Failed to reduce cpu time for pod %q pending volume cleanup due to %v", uid, err)
			}
			continue
		}
		klog.V(3).Infof("Orphaned pod %q found, removing pod cgroups", uid)
		// Destroy all cgroups of pod that should not be running,
		// by first killing all the attached processes to these cgroups.
		// We ignore errors thrown by the method, as the housekeeping loop would
		// again try to delete these unwanted pod cgroups
		go pcm.Destroy(val)
	}
}

// enableHostUserNamespace determines if the host user namespace should be used by the container runtime.
// Returns true if the pod is using a host pid, pic, or network namespace, the pod is using a non-namespaced
// capability, the pod contains a privileged container, or the pod has a host path volume.
//
// NOTE: when if a container shares any namespace with another container it must also share the user namespace
// or it will not have the correct capabilities in the namespace.  This means that host user namespace
// is enabled per pod, not per container.
func (kl *Kubelet) enableHostUserNamespace(pod *v1.Pod) bool {
	if kubecontainer.HasPrivilegedContainer(pod) || hasHostNamespace(pod) ||
		hasHostVolume(pod) || hasNonNamespacedCapability(pod) || kl.hasHostMountPVC(pod) {
		return true
	}
	return false
}

// hasNonNamespacedCapability returns true if MKNOD, SYS_TIME, or SYS_MODULE is requested for any container.
func hasNonNamespacedCapability(pod *v1.Pod) bool {
	for _, c := range pod.Spec.Containers {
		if c.SecurityContext != nil && c.SecurityContext.Capabilities != nil {
			for _, cap := range c.SecurityContext.Capabilities.Add {
				if cap == "MKNOD" || cap == "SYS_TIME" || cap == "SYS_MODULE" {
					return true
				}
			}
		}
	}

	return false
}

// hasHostVolume returns true if the pod spec has a HostPath volume.
func hasHostVolume(pod *v1.Pod) bool {
	for _, v := range pod.Spec.Volumes {
		if v.HostPath != nil {
			return true
		}
	}
	return false
}

// hasHostNamespace returns true if hostIPC, hostNetwork, or hostPID are set to true.
func hasHostNamespace(pod *v1.Pod) bool {
	if pod.Spec.SecurityContext == nil {
		return false
	}
	return pod.Spec.HostIPC || pod.Spec.HostNetwork || pod.Spec.HostPID
}

// hasHostMountPVC returns true if a PVC is referencing a HostPath volume.
func (kl *Kubelet) hasHostMountPVC(pod *v1.Pod) bool {
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pvc, err := kl.kubeClient.CoreV1().PersistentVolumeClaims(pod.Namespace).Get(volume.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
			if err != nil {
				klog.Warningf("unable to retrieve pvc %s:%s - %v", pod.Namespace, volume.PersistentVolumeClaim.ClaimName, err)
				continue
			}
			if pvc != nil {
				referencedVolume, err := kl.kubeClient.CoreV1().PersistentVolumes().Get(pvc.Spec.VolumeName, metav1.GetOptions{})
				if err != nil {
					klog.Warningf("unable to retrieve pv %s - %v", pvc.Spec.VolumeName, err)
					continue
				}
				if referencedVolume != nil && referencedVolume.Spec.HostPath != nil {
					return true
				}
			}
		}
	}
	return false
}
