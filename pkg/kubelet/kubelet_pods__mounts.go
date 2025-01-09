package kubelet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	v1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/features"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/util/format"
	volumeutil "k8s.io/kubernetes/pkg/volume/util"
	"k8s.io/kubernetes/pkg/volume/util/hostutil"
	"k8s.io/kubernetes/pkg/volume/util/subpath"
	volumevalidation "k8s.io/kubernetes/pkg/volume/validation"
)

// makeMounts 挂载 /etc/hosts 文件, 创建 hostPath 类型的 volume 路径.
// 不过返回的是否只是 hostPath 类型的 volume 列表???
//
// 	@param podDir: /var/lib/kubelet/pods/PodUID 
// 	@param podVolumes: 从 actualStateOfWorld 中获取该 Pod 已经声明的 volume 列表(map)
//
// 	@return cleanupAction: Pod移除时, 这些卷的移除方法, 做为回调函数.
//
// caller:
// 	1. pkg/kubelet/kubelet_pods.go -> Kubelet.GenerateRunContainerOptions()
//  只有这一处
//
// makeMounts determines the mount points for the given container.
func makeMounts(
	pod *v1.Pod, podDir string, container *v1.Container, hostName, hostDomain string, 
	podIPs []string, podVolumes kubecontainer.VolumeMap, hu hostutil.HostUtils, 
	subpather subpath.Interface, expandEnvs []kubecontainer.EnvVar,
) ([]kubecontainer.Mount, func(), error) {
	// 是否挂载 /etc/hosts 文件.
	//
	// Kubernetes only mounts on /etc/hosts if:
	// - container is not an infrastructure (pause) container
	// - container is not already mounting on /etc/hosts
	// - OS is not Windows
	// Kubernetes will not mount /etc/hosts if:
	// - when the Pod sandbox is being created, its IP is still unknown.
	// Hence, PodIP will not have been set.
	mountEtcHostsFile := len(podIPs) > 0 && runtime.GOOS != "windows"
	klog.V(3).Infof(
		"container: %v/%v/%v podIPs: %q creating hosts mount: %v", 
		pod.Namespace, pod.Name, container.Name, podIPs, mountEtcHostsFile,
	)
	mounts := []kubecontainer.Mount{}
	var cleanupAction func()
	for i, mount := range container.VolumeMounts {
		// 如果上面已经决定要挂载 /etc/hosts 了, 就忽略 volumeMounts 中的 /etc/hosts.
		//
		// do not mount /etc/hosts if container is already mounting on the path
		mountEtcHostsFile = mountEtcHostsFile && (mount.MountPath != etcHostsPath)
		vol, ok := podVolumes[mount.Name]
		if !ok || vol.Mounter == nil {
			klog.Errorf(
				"Mount cannot be satisfied for container %q, "+
				"because the volume is missing or the volume mounter is nil: %+v", 
				container.Name, mount,
			)
			return nil, cleanupAction, fmt.Errorf(
				"cannot find volume %q to mount into container %q", 
				mount.Name, container.Name,
			)
		}

		relabelVolume := false
		// If the volume supports SELinux and it has not been
		// relabeled already and it is not a read-only volume,
		// relabel it and mark it as labeled
		if vol.Mounter.GetAttributes().Managed && vol.Mounter.GetAttributes().SupportsSELinux && !vol.SELinuxLabeled {
			vol.SELinuxLabeled = true
			relabelVolume = true
		}
		// hostPath `volumes`字段中 hostPath 所指定的宿主机路径.
		hostPath, err := volumeutil.GetPath(vol.Mounter)
		if err != nil {
			return nil, cleanupAction, err
		}

		subPath := mount.SubPath
		// 如果 SubPathExpr 字段不为空, 则 SubPath 路径还需要重新计算.
		if mount.SubPathExpr != "" {
			if !utilfeature.DefaultFeatureGate.Enabled(features.VolumeSubpath) {
				return nil, cleanupAction, fmt.Errorf("volume subpaths are disabled")
			}

			if !utilfeature.DefaultFeatureGate.Enabled(features.VolumeSubpathEnvExpansion) {
				return nil, cleanupAction, fmt.Errorf("volume subpath expansion is disabled")
			}

			subPath, err = kubecontainer.ExpandContainerVolumeMounts(mount, expandEnvs)

			if err != nil {
				return nil, cleanupAction, err
			}
		}

		if subPath != "" {
			if !utilfeature.DefaultFeatureGate.Enabled(features.VolumeSubpath) {
				return nil, cleanupAction, fmt.Errorf("volume subpaths are disabled")
			}

			if filepath.IsAbs(subPath) {
				return nil, cleanupAction, fmt.Errorf(
					"error SubPath `%s` must not be an absolute path", subPath,
				)
			}

			err = volumevalidation.ValidatePathNoBacksteps(subPath)
			if err != nil {
				return nil, cleanupAction, fmt.Errorf(
					"unable to provision SubPath `%s`: %v", subPath, err,
				)
			}

			volumePath := hostPath
			hostPath = filepath.Join(volumePath, subPath)

			if subPathExists, err := hu.PathExists(hostPath); err != nil {
				klog.Errorf(
					"Could not determine if subPath %s exists; "+
					"will not attempt to change its permissions", 
					hostPath,
				)
			} else if !subPathExists {
				// 目标目录不存在则创建.
				//
				// Create the sub path now because if it's auto-created later when referenced,
				// it may have an incorrect ownership and mode. 
				// For example, the sub path directory must have at least g+rwx
				// when the pod specifies an fsGroup, and if the directory is not created here,
				// Docker will later auto-create it with the incorrect mode 0750
				// Make extra care not to escape the volume!
				perm, err := hu.GetMode(volumePath)
				if err != nil {
					return nil, cleanupAction, err
				}
				if err := subpather.SafeMakeDir(subPath, volumePath, perm); err != nil {
					// Don't pass detailed error back to the user 
					// because it could give information about host filesystem
					klog.Errorf(
						"failed to create subPath directory for volumeMount %q of container %q: %v", 
						mount.Name, container.Name, err,
					)
					return nil, cleanupAction, fmt.Errorf(
						"failed to create subPath directory for volumeMount %q of container %q", 
						mount.Name, container.Name,
					)
				}
			}
			hostPath, cleanupAction, err = subpather.PrepareSafeSubpath(subpath.Subpath{
				VolumeMountIndex: i,
				Path:             hostPath,
				VolumeName:       vol.InnerVolumeSpecName,
				VolumePath:       volumePath,
				PodDir:           podDir,
				ContainerName:    container.Name,
			})
			if err != nil {
				// Don't pass detailed error back to the user because it could
				// give information about host filesystem
				klog.Errorf(
					"failed to prepare subPath for volumeMount %q of container %q: %v", 
					mount.Name, container.Name, err,
				)
				return nil, cleanupAction, fmt.Errorf(
					"failed to prepare subPath for volumeMount %q of container %q", 
					mount.Name, container.Name,
				)
			}
		}

		containerPath := mount.MountPath
		propagation, err := translateMountPropagation(mount.MountPropagation)
		if err != nil {
			return nil, cleanupAction, err
		}
		klog.V(5).Infof(
			"Pod %q container %q mount %q has propagation %q", 
			format.Pod(pod), container.Name, mount.Name, propagation,
		)

		mustMountRO := vol.Mounter.GetAttributes().ReadOnly

		mounts = append(mounts, kubecontainer.Mount{
			Name:           mount.Name,
			ContainerPath:  containerPath,
			HostPath:       hostPath,
			ReadOnly:       mount.ReadOnly || mustMountRO,
			SELinuxRelabel: relabelVolume,
			Propagation:    propagation,
		})
	}
	if mountEtcHostsFile {
		hostAliases := pod.Spec.HostAliases
		hostsMount, err := makeHostsMount(
			podDir, podIPs, hostName, hostDomain, hostAliases, pod.Spec.HostNetwork,
		)
		if err != nil {
			return nil, cleanupAction, err
		}
		mounts = append(mounts, *hostsMount)
	}
	return mounts, cleanupAction, nil
}

// translateMountPropagation transforms v1.MountPropagationMode to
// runtimeapi.MountPropagation.
func translateMountPropagation(mountMode *v1.MountPropagationMode) (runtimeapi.MountPropagation, error) {
	switch {
	case mountMode == nil:
		// PRIVATE is the default
		return runtimeapi.MountPropagation_PROPAGATION_PRIVATE, nil
	case *mountMode == v1.MountPropagationHostToContainer:
		return runtimeapi.MountPropagation_PROPAGATION_HOST_TO_CONTAINER, nil
	case *mountMode == v1.MountPropagationBidirectional:
		return runtimeapi.MountPropagation_PROPAGATION_BIDIRECTIONAL, nil
	case *mountMode == v1.MountPropagationNone:
		return runtimeapi.MountPropagation_PROPAGATION_PRIVATE, nil
	default:
		return 0, fmt.Errorf("invalid MountPropagation mode: %q", *mountMode)
	}
}

// makeHostsMount 在宿主机为目标 Pod 准备 hosts 文件, 之后会将其挂载到 /etc/hosts 
// 会根据目标 Pod 是否为 hostNetwork 写入相应的配置.
//
// 	@param podDir: /var/lib/kubelet/pods/${podUID}
//
// makeHostsMount makes the mountpoint for the hosts file that the containers
// in a pod are injected with. podIPs is provided instead of podIP as podIPs
// are present even if dual-stack feature flag is not enabled.
func makeHostsMount(
	podDir string, podIPs []string, 
	hostName, hostDomainName string, hostAliases []v1.HostAlias, 
	useHostNetwork bool,
) (*kubecontainer.Mount, error) {
	// /var/lib/kubelet/pods/${podUID}/etc-hosts
	hostsFilePath := path.Join(podDir, "etc-hosts")
	err := ensureHostsFile(
		hostsFilePath, podIPs, hostName, hostDomainName, hostAliases, 
		useHostNetwork,
	)
	if err != nil {
		return nil, err
	}
	return &kubecontainer.Mount{
		Name:           "k8s-managed-etc-hosts",
		ContainerPath:  etcHostsPath,
		HostPath:       hostsFilePath,
		ReadOnly:       false,
		SELinuxRelabel: true,
	}, nil
}

// ensureHostsFile ensures that the given host file has an up-to-date ip, host
// name, and domain name.
func ensureHostsFile(
	fileName string, hostIPs []string, 
	hostName, hostDomainName string, hostAliases []v1.HostAlias, 
	useHostNetwork bool,
) error {
	var hostsFileContent []byte
	var err error

	if useHostNetwork {
		// if Pod is using host network, read hosts file from the node's filesystem.
		// `etcHostsPath` references the location of the hosts file on the node.
		// `/etc/hosts` for *nix systems.
		hostsFileContent, err = nodeHostsFileContent(etcHostsPath, hostAliases)
		if err != nil {
			return err
		}
	} else {
		// if Pod is not using host network, create a managed hosts file with Pod IP and other information.
		hostsFileContent = managedHostsFileContent(hostIPs, hostName, hostDomainName, hostAliases)
	}

	return ioutil.WriteFile(fileName, hostsFileContent, 0644)
}

// nodeHostsFileContent reads the content of node's hosts file.
func nodeHostsFileContent(hostsFilePath string, hostAliases []v1.HostAlias) ([]byte, error) {
	hostsFileContent, err := ioutil.ReadFile(hostsFilePath)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString(managedHostsHeaderWithHostNetwork)
	buffer.Write(hostsFileContent)
	buffer.Write(hostsEntriesFromHostAliases(hostAliases))
	return buffer.Bytes(), nil
}

// managedHostsFileContent generates the content of the managed etc hosts based on Pod IPs and other
// information.
func managedHostsFileContent(hostIPs []string, hostName, hostDomainName string, hostAliases []v1.HostAlias) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(managedHostsHeader)
	buffer.WriteString("127.0.0.1\tlocalhost\n")                      // ipv4 localhost
	buffer.WriteString("::1\tlocalhost ip6-localhost ip6-loopback\n") // ipv6 localhost
	buffer.WriteString("fe00::0\tip6-localnet\n")
	buffer.WriteString("fe00::0\tip6-mcastprefix\n")
	buffer.WriteString("fe00::1\tip6-allnodes\n")
	buffer.WriteString("fe00::2\tip6-allrouters\n")
	if len(hostDomainName) > 0 {
		// host entry generated for all IPs in podIPs
		// podIPs field is populated for clusters even
		// dual-stack feature flag is not enabled.
		for _, hostIP := range hostIPs {
			buffer.WriteString(fmt.Sprintf("%s\t%s.%s\t%s\n", hostIP, hostName, hostDomainName, hostName))
		}
	} else {
		for _, hostIP := range hostIPs {
			buffer.WriteString(fmt.Sprintf("%s\t%s\n", hostIP, hostName))
		}
	}
	buffer.Write(hostsEntriesFromHostAliases(hostAliases))
	return buffer.Bytes()
}

func hostsEntriesFromHostAliases(hostAliases []v1.HostAlias) []byte {
	if len(hostAliases) == 0 {
		return []byte{}
	}

	var buffer bytes.Buffer
	buffer.WriteString("\n")
	buffer.WriteString("# Entries added by HostAliases.\n")
	// for each IP, write all aliases onto single line in hosts file
	for _, hostAlias := range hostAliases {
		buffer.WriteString(fmt.Sprintf("%s\t%s\n", hostAlias.IP, strings.Join(hostAlias.Hostnames, "\t")))
	}
	return buffer.Bytes()
}

// truncatePodHostnameIfNeeded truncates the pod hostname if it's longer than 63 chars.
func truncatePodHostnameIfNeeded(podName, hostname string) (string, error) {
	// Cap hostname at 63 chars (specification is 64bytes which is 63 chars and the null terminating char).
	const hostnameMaxLen = 63
	if len(hostname) <= hostnameMaxLen {
		return hostname, nil
	}
	truncated := hostname[:hostnameMaxLen]
	klog.Errorf("hostname for pod:%q was longer than %d. Truncated hostname to :%q", podName, hostnameMaxLen, truncated)
	// hostname should not end with '-' or '.'
	truncated = strings.TrimRight(truncated, "-.")
	if len(truncated) == 0 {
		// This should never happen.
		return "", fmt.Errorf("hostname for pod %q was invalid: %q", podName, hostname)
	}
	return truncated, nil
}
