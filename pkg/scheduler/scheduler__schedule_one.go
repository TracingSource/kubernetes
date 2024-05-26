package scheduler

import (
	"fmt"
	"time"
	"context"
	"math/rand"

	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/core"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/metrics"
)

// scheduleOne 调度单个 Pod 经历的全部流程
//
// caller:
// 	1. cmd/kube-scheduler/app/server.go -> Run() 一般在选主后由主节点被循环调用.
//
// scheduleOne does the entire scheduling workflow for a single pod. 
// It is serialized on the scheduling algorithm's host fitting.
func (sched *Scheduler) scheduleOne(ctx context.Context) {
	fwk := sched.Framework

	// 从待调度队列(SchedulingQueue 成员变量)中取出一个 Pod(如没有则先阻塞).
	podInfo := sched.NextPod()
	// pod could be nil when schedulerQueue is closed
	if podInfo == nil || podInfo.Pod == nil {
		return
	}
	pod := podInfo.Pod
	if pod.DeletionTimestamp != nil {
		sched.Recorder.Eventf(
			pod, nil, v1.EventTypeWarning, "FailedScheduling", "Scheduling", 
			"skip schedule deleting pod: %v/%v", pod.Namespace, pod.Name,
		)
		klog.V(3).Infof("Skip schedule deleting pod: %v/%v", pod.Namespace, pod.Name)
		return
	}

	klog.V(3).Infof("Attempting to schedule pod: %v/%v", pod.Namespace, pod.Name)

	// Synchronously attempt to find a fit for the pod.
	start := time.Now()
	state := framework.NewCycleState()
	state.SetRecordFrameworkMetrics(rand.Intn(100) < frameworkMetricsSamplePercent)
	schedulingCycleCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// 这里运行 "预选", "优选" 两个阶段的调度算法.
	scheduleResult, err := sched.Algorithm.Schedule(schedulingCycleCtx, state, pod)
	if err != nil {
		// Pod 调度失败时, 发送一条调度失败的 event, 并将失败原因更新到 pod.Status.Conditions 列表中.
		sched.recordSchedulingFailure(podInfo.DeepCopy(), err, v1.PodReasonUnschedulable, err.Error())
		// 调度过程可能失败, 目标 pod 找不到合适的 host.
		// 这样只能尝试 preempt(抢占), 希望抢占能够调度上去.
		// 当然, 也有可能被其他 pod 调上去, 抢占不到, 这种情况是正常的.
		//
		// Schedule() may have failed because the pod would not fit on any host,
		// so we try to preempt, with the expectation that the next time the pod
		// is tried for scheduling it will fit due to the preemption.
		// It is also possible that a different pod will schedule
		// into the resources that were preempted, but this is harmless.
		if fitError, ok := err.(*core.FitError); ok {
			if sched.DisablePreemption {
				klog.V(3).Infof(
					"Pod priority feature is not enabled or preemption is disabled by scheduler configuration." +
					" No preemption is performed.",
				)
			} else {
				preemptionStartTime := time.Now()
				// 抢占
				sched.preempt(schedulingCycleCtx, state, fwk, pod, fitError)
				metrics.PreemptionAttempts.Inc()
				metrics.SchedulingAlgorithmPreemptionEvaluationDuration.Observe(metrics.SinceInSeconds(preemptionStartTime))
				metrics.DeprecatedSchedulingAlgorithmPreemptionEvaluationDuration.Observe(metrics.SinceInMicroseconds(preemptionStartTime))
				metrics.SchedulingLatency.WithLabelValues(metrics.PreemptionEvaluation).Observe(metrics.SinceInSeconds(preemptionStartTime))
				metrics.DeprecatedSchedulingLatency.WithLabelValues(metrics.PreemptionEvaluation).Observe(metrics.SinceInSeconds(preemptionStartTime))
			}
			// Pod did not fit anywhere, so it is counted as a failure.
			// If preemption succeeds, the pod should get counted as a success
			// the next time we try to schedule it. (hopefully)
			metrics.PodScheduleFailures.Inc()
		} else {
			klog.Errorf("error selecting node for pod: %v", err)
			metrics.PodScheduleErrors.Inc()
		}
		return
	}
	metrics.SchedulingAlgorithmLatency.Observe(metrics.SinceInSeconds(start))
	metrics.DeprecatedSchedulingAlgorithmLatency.Observe(metrics.SinceInMicroseconds(start))
	// Tell the cache to assume that a pod now is running on a given node,
	// even though it hasn't been bound yet.
	// This allows us to keep scheduling without waiting on binding to occur.
	assumedPodInfo := podInfo.DeepCopy()
	assumedPod := assumedPodInfo.Pod

	// Assume volumes first before assuming the pod.
	//
	// If all volumes are completely bound, then allBound is true and binding will be skipped.
	//
	// Otherwise, binding of volumes is started after the pod is assumed, but before pod binding.
	//
	// This function modifies 'assumedPod' if volume binding is required.
	allBound, err := sched.VolumeBinder.Binder.AssumePodVolumes(
		assumedPod, scheduleResult.SuggestedHost,
	)
	if err != nil {
		sched.recordSchedulingFailure(
			assumedPodInfo, err, SchedulerError,
			fmt.Sprintf("AssumePodVolumes failed: %v", err),
		)
		metrics.PodScheduleErrors.Inc()
		return
	}

	// Run "reserve" plugins.
	sts := fwk.RunReservePlugins(
		schedulingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost,
	)
	if !sts.IsSuccess() {
		sched.recordSchedulingFailure(
			assumedPodInfo, sts.AsError(), SchedulerError, sts.Message(),
		)
		metrics.PodScheduleErrors.Inc()
		return
	}

	// assume modifies `assumedPod` by setting NodeName=scheduleResult.SuggestedHost
	err = sched.assume(assumedPod, scheduleResult.SuggestedHost)
	if err != nil {
		// This is most probably result of a BUG in retrying logic.
		// We report an error here so that pod scheduling can be retried.
		// This relies on the fact that Error will check if the pod has been bound
		// to a node and if so will not add it back to the unscheduled pods queue
		// (otherwise this would cause an infinite loop).
		sched.recordSchedulingFailure(
			assumedPodInfo, err, SchedulerError, fmt.Sprintf("AssumePod failed: %v", err),
		)
		metrics.PodScheduleErrors.Inc()
		// trigger un-reserve plugins to clean up state associated with the reserved Pod
		fwk.RunUnreservePlugins(schedulingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
		return
	}
	// bind the pod to its host asynchronously (we can do this b/c of the assumption step above).
	go func() {
		bindingCycleCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		metrics.SchedulerGoroutines.WithLabelValues("binding").Inc()
		defer metrics.SchedulerGoroutines.WithLabelValues("binding").Dec()

		// Run "permit" plugins.
		permitStatus := fwk.RunPermitPlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
		if !permitStatus.IsSuccess() {
			var reason string
			if permitStatus.IsUnschedulable() {
				metrics.PodScheduleFailures.Inc()
				reason = v1.PodReasonUnschedulable
			} else {
				metrics.PodScheduleErrors.Inc()
				reason = SchedulerError
			}
			if forgetErr := sched.Cache().ForgetPod(assumedPod); forgetErr != nil {
				klog.Errorf("scheduler cache ForgetPod failed: %v", forgetErr)
			}
			// trigger un-reserve plugins to clean up state associated with the reserved Pod
			fwk.RunUnreservePlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
			sched.recordSchedulingFailure(assumedPodInfo, permitStatus.AsError(), reason, permitStatus.Message())
			return
		}

		// Bind volumes first before Pod
		if !allBound {
			err := sched.bindVolumes(assumedPod)
			if err != nil {
				sched.recordSchedulingFailure(assumedPodInfo, err, "VolumeBindingFailed", err.Error())
				metrics.PodScheduleErrors.Inc()
				// trigger un-reserve plugins to clean up state associated with the reserved Pod
				fwk.RunUnreservePlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
				return
			}
		}

		// Run "prebind" plugins.
		preBindStatus := fwk.RunPreBindPlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
		if !preBindStatus.IsSuccess() {
			var reason string
			metrics.PodScheduleErrors.Inc()
			reason = SchedulerError
			if forgetErr := sched.Cache().ForgetPod(assumedPod); forgetErr != nil {
				klog.Errorf("scheduler cache ForgetPod failed: %v", forgetErr)
			}
			// trigger un-reserve plugins to clean up state associated with the reserved Pod
			fwk.RunUnreservePlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
			sched.recordSchedulingFailure(assumedPodInfo, preBindStatus.AsError(), reason, preBindStatus.Message())
			return
		}

		err := sched.bind(bindingCycleCtx, assumedPod, scheduleResult.SuggestedHost, state)
		metrics.E2eSchedulingLatency.Observe(metrics.SinceInSeconds(start))
		metrics.DeprecatedE2eSchedulingLatency.Observe(metrics.SinceInMicroseconds(start))
		if err != nil {
			metrics.PodScheduleErrors.Inc()
			// trigger un-reserve plugins to clean up state associated with the reserved Pod
			fwk.RunUnreservePlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
			sched.recordSchedulingFailure(assumedPodInfo, err, SchedulerError, fmt.Sprintf("Binding rejected: %v", err))
		} else {
			// Calculating nodeResourceString can be heavy. Avoid it if klog verbosity is below 2.
			if klog.V(2) {
				klog.Infof("pod %v/%v is bound successfully on node %q, %d nodes evaluated, %d nodes were found feasible.", assumedPod.Namespace, assumedPod.Name, scheduleResult.SuggestedHost, scheduleResult.EvaluatedNodes, scheduleResult.FeasibleNodes)
			}

			metrics.PodScheduleSuccesses.Inc()
			metrics.PodSchedulingAttempts.Observe(float64(podInfo.Attempts))
			metrics.PodSchedulingDuration.Observe(metrics.SinceInSeconds(podInfo.InitialAttemptTimestamp))

			// Run "postbind" plugins.
			fwk.RunPostBindPlugins(bindingCycleCtx, state, assumedPod, scheduleResult.SuggestedHost)
		}
	}()
}
