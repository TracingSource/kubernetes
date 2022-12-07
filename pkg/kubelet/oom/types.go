package oom

import v1 "k8s.io/api/core/v1"

// Watcher 由 pkg/kubelet/oom/oom_watcher_linux.go -> realWatcher{} 结构体实现
//
// Watcher defines the interface of OOM watchers.
type Watcher interface {
	Start(ref *v1.ObjectReference) error
}
