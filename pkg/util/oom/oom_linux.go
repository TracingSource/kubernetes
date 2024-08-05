//go:build linux
// +build linux

package oom

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	cmutil "k8s.io/kubernetes/pkg/kubelet/cm/util"

	"k8s.io/klog"
)

// caller:
// 	1. cmd/kubelet/app/server.go -> UnsecuredDependencies() 在 kubelet 启动过程中被调用.
// 	2. cmd/kube-proxy/app/server__proxy_server.go -> ProxyServer.Run()
// 	3. pkg/kubelet/cm/container_manager_linux.go -> ensureProcessInContainerWithOOMScore()
func NewOOMAdjuster() *OOMAdjuster {
	oomAdjuster := &OOMAdjuster{
		pidLister:        getPids,
		ApplyOOMScoreAdj: applyOOMScoreAdj,
	}
	oomAdjuster.ApplyOOMScoreAdjContainer = oomAdjuster.applyOOMScoreAdjContainer
	return oomAdjuster
}

func getPids(cgroupName string) ([]int, error) {
	return cmutil.GetPids(filepath.Join("/", cgroupName))
}

// applyOOMScoreAdj 为目标 pid 进程设置 oom_score_adj 值
//
// 	@param pid: 目标进程(如果为0则表示为自身 self)
// 	@param oomScoreAdj: oom_score_adj 值
//
// caller:
// 	1. cmd/kube-proxy/app/server__proxy_server.go -> ProxyServer.Run()
// 	2. cmd/kubelet/app/server.go -> run()
// 	3. pkg/kubelet/cm/container_manager_linux.go -> ensureProcessInContainerWithOOMScore()
// 	4. OOMAdjuster.applyOOMScoreAdjContainer()
//
// Writes 'value' to /proc/<pid>/oom_score_adj. PID = 0 means self
// Returns os.ErrNotExist if the `pid` does not exist.
func applyOOMScoreAdj(pid int, oomScoreAdj int) error {
	if pid < 0 {
		return fmt.Errorf("invalid PID %d specified for oom_score_adj", pid)
	}

	var pidStr string
	if pid == 0 {
		pidStr = "self"
	} else {
		pidStr = strconv.Itoa(pid)
	}

	maxTries := 2
	oomScoreAdjPath := path.Join("/proc", pidStr, "oom_score_adj")
	value := strconv.Itoa(oomScoreAdj)
	klog.V(4).Infof("attempting to set %q to %q", oomScoreAdjPath, value)
	var err error
	for i := 0; i < maxTries; i++ {
		err = ioutil.WriteFile(oomScoreAdjPath, []byte(value), 0700)
		if err != nil {
			if os.IsNotExist(err) {
				klog.V(2).Infof("%q does not exist", oomScoreAdjPath)
				return os.ErrNotExist
			}

			klog.V(3).Info(err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return nil
	}
	if err != nil {
		klog.V(2).Infof("failed to set %q to %q: %v", oomScoreAdjPath, value, err)
	}
	return err
}

// applyOOMScoreAdjContainer  为目标 cgroup 控制组下所有进程都设置 oom_score_adj 值.
//
// Writes 'value' to /proc/<pid>/oom_score_adj for all processes in cgroup cgroupName.
// Keeps trying to write until the process list of the cgroup stabilizes,
// or until maxTries tries.
func (oomAdjuster *OOMAdjuster) applyOOMScoreAdjContainer(
	cgroupName string, oomScoreAdj, maxTries int,
) error {
	adjustedProcessSet := make(map[int]bool)
	for i := 0; i < maxTries; i++ {
		continueAdjusting := false
		pidList, err := oomAdjuster.pidLister(cgroupName)
		if err != nil {
			if os.IsNotExist(err) {
				// Nothing to do since the container doesn't exist anymore.
				return os.ErrNotExist
			}
			continueAdjusting = true
			klog.V(10).Infof("Error getting process list for cgroup %s: %+v", cgroupName, err)
		} else if len(pidList) == 0 {
			klog.V(10).Infof("Pid list is empty")
			continueAdjusting = true
		} else {
			for _, pid := range pidList {
				if !adjustedProcessSet[pid] {
					klog.V(10).Infof("pid %d needs to be set", pid)
					if err = oomAdjuster.ApplyOOMScoreAdj(pid, oomScoreAdj); err == nil {
						adjustedProcessSet[pid] = true
					} else if err == os.ErrNotExist {
						continue
					} else {
						klog.V(10).Infof("cannot adjust oom score for pid %d - %v", pid, err)
						continueAdjusting = true
					}
					// Processes can come and go while we try to apply oom score adjust value.
					// So ignore errors here.
				}
			}
		}
		if !continueAdjusting {
			return nil
		}
		// There's a slight race. 
		// A process might have forked just before we write its OOM score adjust.
		// The fork might copy the parent process's old OOM score,
		// then this function might execute and update the parent's OOM score,
		// but the forked process id might not be reflected in cgroup.procs
		// for a short amount of time.
		// So this function might return without changing the forked process's
		// OOM score.
		// Very unlikely race, so ignoring this for now.
	}
	return fmt.Errorf("exceeded maxTries, some processes might not have desired OOM score")
}
