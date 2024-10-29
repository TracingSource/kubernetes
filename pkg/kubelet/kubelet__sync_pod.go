package kubelet

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"k8s.io/kubernetes/pkg/kubelet/images"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/kubelet/util/format"
)

// caller:
// 	1. pkg/kubelet/pod_workers.go -> managePodLoop() 作为 syncPodFn() 被调用.
//
// syncPod is the transaction script for the sync of a single pod.
//
// Arguments:
//
// o - the SyncPodOptions for this invocation
//
// The workflow is:
// * If the pod is being created, record pod worker start latency
// * Call generateAPIPodStatus to prepare an v1.PodStatus for the pod
// * If the pod is being seen as running for the first time, record pod
//   start latency
// * Update the status of the pod in the status manager
// * Kill the pod if it should not be running
// * Create a mirror pod if the pod is a static pod, and does not
//   already have a mirror pod
// * Create the data directories for the pod if they do not exist
// * Wait for volumes to attach/mount
// * Fetch the pull secrets for the pod
// * Call the container runtime's SyncPod callback
// * Update the traffic shaping for the pod's ingress and egress limits
//
// If any step of this workflow errors, the error is returned, and is repeated
// on the next syncPod call.
//
// This operation writes all events that are dispatched in order to provide
// the most accurate information possible about an error situation to aid debugging.
// Callers should not throw an event if this operation returns an error.
func (kl *Kubelet) syncPod(o syncPodOptions) error {
	// pull out the required options
	pod := o.pod
	mirrorPod := o.mirrorPod
	podStatus := o.podStatus
	updateType := o.updateType

	// if we want to kill a pod, do it now!
	if updateType == kubetypes.SyncPodKill {
		killPodOptions := o.killPodOptions
		if killPodOptions == nil || killPodOptions.PodStatusFunc == nil {
			return fmt.Errorf("kill pod options are required if update type is kill")
		}
		apiPodStatus := killPodOptions.PodStatusFunc(pod, podStatus)
		kl.statusManager.SetPodStatus(pod, apiPodStatus)
		// we kill the pod with the specified grace period since this is a termination
		if err := kl.killPod(pod, nil, podStatus, killPodOptions.PodTerminationGracePeriodSecondsOverride); err != nil {
			kl.recorder.Eventf(pod, v1.EventTypeWarning, events.FailedToKillPod, "error killing pod: %v", err)
			// there was an error killing the pod, so we return that error directly
			utilruntime.HandleError(err)
			return err
		}
		return nil
	}

	// Latency measurements for the main workflow are relative to the
	// first time the pod was seen by the API server.
	var firstSeenTime time.Time
	if firstSeenTimeStr, ok := pod.Annotations[kubetypes.ConfigFirstSeenAnnotationKey]; ok {
		firstSeenTime = kubetypes.ConvertToTimestamp(firstSeenTimeStr).Get()
	}

	// Record pod worker start latency if being created
	// TODO: make pod workers record their own latencies
	if updateType == kubetypes.SyncPodCreate {
		if !firstSeenTime.IsZero() {
			// This is the first time we are syncing the pod. Record the latency
			// since kubelet first saw the pod if firstSeenTime is set.
			metrics.PodWorkerStartDuration.Observe(metrics.SinceInSeconds(firstSeenTime))
			metrics.DeprecatedPodWorkerStartLatency.Observe(metrics.SinceInMicroseconds(firstSeenTime))
		} else {
			klog.V(3).Infof("First seen time not recorded for pod %q", pod.UID)
		}
	}
	////////////////////////////////////////////////////////////////////////////
	// pod status 相关
	//
	// Generate final API pod status with pod and status manager status
	apiPodStatus := kl.generateAPIPodStatus(pod, podStatus)
	// The pod IP may be changed in generateAPIPodStatus
	// if the pod is using host network. (See #24576)
	// TODO(random-liu): After writing pod spec into container labels,
	// check whether pod is using host network,
	// and set pod IP to hostIP directly in runtime.GetPodStatus
	podStatus.IPs = make([]string, 0, len(apiPodStatus.PodIPs))
	for _, ipInfo := range apiPodStatus.PodIPs {
		podStatus.IPs = append(podStatus.IPs, ipInfo.IP)
	}

	if len(podStatus.IPs) == 0 && len(apiPodStatus.PodIP) > 0 {
		podStatus.IPs = []string{apiPodStatus.PodIP}
	}

	// Record the time it takes for the pod to become running.
	existingStatus, ok := kl.statusManager.GetPodStatus(pod.UID)
	if !ok || existingStatus.Phase == v1.PodPending && apiPodStatus.Phase == v1.PodRunning &&
		!firstSeenTime.IsZero() {
		metrics.PodStartDuration.Observe(metrics.SinceInSeconds(firstSeenTime))
		metrics.DeprecatedPodStartLatency.Observe(metrics.SinceInMicroseconds(firstSeenTime))
	}

	runnable := kl.canRunPod(pod)
	if !runnable.Admit {
		// Pod is not runnable; update the Pod and Container statuses to why.
		apiPodStatus.Reason = runnable.Reason
		apiPodStatus.Message = runnable.Message
		// Waiting containers are not creating.
		const waitingReason = "Blocked"
		for _, cs := range apiPodStatus.InitContainerStatuses {
			if cs.State.Waiting != nil {
				cs.State.Waiting.Reason = waitingReason
			}
		}
		for _, cs := range apiPodStatus.ContainerStatuses {
			if cs.State.Waiting != nil {
				cs.State.Waiting.Reason = waitingReason
			}
		}
	}

	// Update status in the status manager
	kl.statusManager.SetPodStatus(pod, apiPodStatus)
	//
	// pod status 结束
	////////////////////////////////////////////////////////////////////////////

	// Kill pod if it should not be running
	if !runnable.Admit || pod.DeletionTimestamp != nil || apiPodStatus.Phase == v1.PodFailed {
		var syncErr error
		if err := kl.killPod(pod, nil, podStatus, nil); err != nil {
			kl.recorder.Eventf(pod, v1.EventTypeWarning, events.FailedToKillPod, "error killing pod: %v", err)
			syncErr = fmt.Errorf("error killing pod: %v", err)
			utilruntime.HandleError(syncErr)
		} else {
			if !runnable.Admit {
				// There was no error killing the pod, but the pod cannot be run.
				// Return an error to signal that the sync loop should back off.
				syncErr = fmt.Errorf("pod cannot be run: %s", runnable.Message)
			}
		}
		return syncErr
	}

	// 如果网络插件还未就绪, 且目标pod并不是 host network 模式, 则返回错误.
	//
	// If the network plugin is not ready, only start the pod if it uses the host network
	if err := kl.runtimeState.networkErrors(); err != nil && !kubecontainer.IsHostNetworkPod(pod) {
		kl.recorder.Eventf(
			pod, v1.EventTypeWarning, events.NetworkNotReady, 
			"%s: %v", NetworkNotReadyErrorMsg, err,
		)
		return fmt.Errorf("%s: %v", NetworkNotReadyErrorMsg, err)
	}

	// Create Cgroups for the pod and apply resource parameters
	// to them if cgroups-per-qos flag is enabled.
	pcm := kl.containerManager.NewPodContainerManager()
	// If pod has already been terminated then we need not create
	// or update the pod's cgroup
	if !kl.podIsTerminated(pod) {
		// 当 kubelet 之前没有开启 cgroups-per-qos, 后来又开启并重启后,
		// 需要当原本已经在运行的 pod 删除掉, 然后按照 cgroup 层级重新建立.
		// 但是对于第1次建立的 pod 则不需要.
		//
		// When the kubelet is restarted with the cgroups-per-qos
		// flag enabled, all the pod's running containers
		// should be killed intermittently and brought back up
		// under the qos cgroup hierarchy.
		// Check if this is the pod's first sync
		firstSync := true
		for _, containerStatus := range apiPodStatus.ContainerStatuses {
			if containerStatus.State.Running != nil {
				firstSync = false
				break
			}
		}
		// Don't kill containers in pod if pod's cgroups already
		// exists or the pod is running for the first time
		podKilled := false
		// 判断目标 pod 是否拥有**所有**的 cgroup 子系统目录.
		if !pcm.Exists(pod) && !firstSync {
			if err := kl.killPod(pod, nil, podStatus, nil); err == nil {
				podKilled = true
			}
		}
		// 如果 pod 是一次性运行(重启策略为 never), 就不要为ta创建 cgroup 了.
		// ...呃, 那这不是相当于重启一次 kubelet 就把那些 pod 给终止了吗???
		//
		// Create and Update pod's Cgroups
		// Don't create cgroups for run once pod if it was killed above
		// The current policy is not to restart the run once pods when
		// the kubelet is restarted with the new flag as run once pods are
		// expected to run only once and if the kubelet is restarted then
		// they are not expected to run again.
		// We don't create and apply updates to cgroup if its a run once pod and was killed above
		if !(podKilled && pod.Spec.RestartPolicy == v1.RestartPolicyNever) {
			if !pcm.Exists(pod) {
				// 实际调用的是 pkg/kubelet/cm/qos_container_manager_linux.go -> UpdateCgroups()
				if err := kl.containerManager.UpdateQOSCgroups(); err != nil {
					klog.V(2).Infof("Failed to update QoS cgroups while syncing pod: %v", err)
				}
				if err := pcm.EnsureExists(pod); err != nil {
					kl.recorder.Eventf(
						pod, v1.EventTypeWarning, events.FailedToCreatePodContainer, 
						"unable to ensure pod container exists: %v", err,
					)
					return fmt.Errorf("failed to ensure that the pod: %v cgroups exist and are correctly applied: %v", pod.UID, err)
				}
			}
		}
	}

	// Create Mirror Pod for Static Pod if it doesn't already exist
	if kubetypes.IsStaticPod(pod) {
		podFullName := kubecontainer.GetPodFullName(pod)
		deleted := false
		// kubelet 可以自行管理 static pod, 但是还是会向 apiserver 注册对应的 mirror pod,
		// 这也是为何可以通过 kubectl 查看到这些 static pod 的原因.
		//
		// 如果 mirrorPod != nil, 说明该 pod 在 apiserver 中存在一个对应的 mirror pod.
		if mirrorPod != nil {
			// 如果 mirror pod 已经被删除, 或不再符合当前 static pod, 则删除 mirror pod.
			// 注意!!! 删除的是 mirror pod, 因为 mirror pod 是根据 static pod 创建的映射对象.
			if mirrorPod.DeletionTimestamp != nil || !kl.podManager.IsMirrorPodOf(mirrorPod, pod) {
				// The mirror pod is semantically different from the static pod. Remove
				// it. The mirror pod will get recreated later.
				klog.Infof("Trying to delete pod %s %v", podFullName, mirrorPod.ObjectMeta.UID)
				var err error
				deleted, err = kl.podManager.DeleteMirrorPod(podFullName, &mirrorPod.ObjectMeta.UID)
				if deleted {
					klog.Warningf("Deleted mirror pod %q because it is outdated", format.Pod(mirrorPod))
				} else if err != nil {
					klog.Errorf("Failed deleting mirror pod %q: %v", format.Pod(mirrorPod), err)
				}
			}
		}
		if mirrorPod == nil || deleted {
			node, err := kl.GetNode()
			if err != nil || node.DeletionTimestamp != nil {
				klog.V(4).Infof("No need to create a mirror pod, since node %q has been removed from the cluster", kl.nodeName)
			} else {
				klog.V(4).Infof("Creating a mirror pod for static pod %q", format.Pod(pod))
				if err := kl.podManager.CreateMirrorPod(pod); err != nil {
					klog.Errorf("Failed creating a mirror pod for %q: %v", format.Pod(pod), err)
				}
			}
		}
	}

	// Make data directories for the pod
	if err := kl.makePodDataDirs(pod); err != nil {
		kl.recorder.Eventf(
			pod, v1.EventTypeWarning, events.FailedToMakePodDataDirectories, 
			"error making pod data directories: %v", err,
		)
		klog.Errorf("Unable to make pod data directories for pod %q: %v", format.Pod(pod), err)
		return err
	}

	// Volume manager will not mount volumes for terminated pods
	if !kl.podIsTerminated(pod) {
		// Wait for volumes to attach/mount
		if err := kl.volumeManager.WaitForAttachAndMount(pod); err != nil {
			kl.recorder.Eventf(
				pod, v1.EventTypeWarning, events.FailedMountVolume, 
				"Unable to attach or mount volumes: %v", err,
			)
			klog.Errorf(
				"Unable to attach or mount volumes for pod %q: %v; skipping pod", 
				format.Pod(pod), err,
			)
			return err
		}
	}

	// Fetch the pull secrets for the pod
	pullSecrets := kl.getPullSecretsForPod(pod)

	// Call the container runtime's SyncPod callback
	result := kl.containerRuntime.SyncPod(pod, podStatus, pullSecrets, kl.backOff)
	kl.reasonCache.Update(pod.UID, result)
	if err := result.Error(); err != nil {
		// Do not return error if the only failures were pods in backoff
		for _, r := range result.SyncResults {
			if r.Error != kubecontainer.ErrCrashLoopBackOff && r.Error != images.ErrImagePullBackOff {
				// Do not record an event here, as we keep all event logging for sync pod failures
				// local to container runtime so we get better errors
				return err
			}
		}

		return nil
	}

	return nil
}
