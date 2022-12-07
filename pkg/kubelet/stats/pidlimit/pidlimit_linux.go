// +build linux

package pidlimit

import (
	"io/ioutil"
	"strconv"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// Stats provides basic information about max and current process count
func Stats() (*statsapi.RlimitStats, error) {
	rlimit := &statsapi.RlimitStats{}

	if content, err := ioutil.ReadFile("/proc/sys/kernel/pid_max"); err == nil {
		if maxPid, err := strconv.ParseInt(string(content[:len(content)-1]), 10, 64); err == nil {
			rlimit.MaxPID = &maxPid
		}
	}

	var info syscall.Sysinfo_t
	syscall.Sysinfo(&info)
	procs := int64(info.Procs)
	rlimit.NumOfRunningProcesses = &procs

	rlimit.Time = v1.NewTime(time.Now())

	return rlimit, nil
}
