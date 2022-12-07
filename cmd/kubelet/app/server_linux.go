package app

import (
	"k8s.io/klog"
	"k8s.io/utils/inotify"
)

func watchForLockfileContention(path string, done chan struct{}) error {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		klog.Errorf("unable to create watcher for lockfile: %v", err)
		return err
	}
	if err = watcher.AddWatch(path, inotify.InOpen|inotify.InDeleteSelf); err != nil {
		klog.Errorf("unable to watch lockfile: %v", err)
		return err
	}
	go func() {
		select {
		case ev := <-watcher.Event:
			klog.Infof("inotify event: %v", ev)
		case err = <-watcher.Error:
			klog.Errorf("inotify watcher error: %v", err)
		}
		close(done)
	}()
	return nil
}
