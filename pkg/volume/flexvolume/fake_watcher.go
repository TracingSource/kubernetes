package flexvolume

import (
	"github.com/fsnotify/fsnotify"
	utilfs "k8s.io/kubernetes/pkg/util/filesystem"
)

// Mock filesystem watcher
type fakeWatcher struct {
	watches      []string // List of watches added by the prober, ordered from least recent to most recent.
	eventHandler utilfs.FSEventHandler
}

var _ utilfs.FSWatcher = &fakeWatcher{}

func newFakeWatcher() *fakeWatcher {
	return &fakeWatcher{
		watches: nil,
	}
}

func (w *fakeWatcher) Init(eventHandler utilfs.FSEventHandler, _ utilfs.FSErrorHandler) error {
	w.eventHandler = eventHandler
	return nil
}

func (w *fakeWatcher) Run() { /* no-op */ }

func (w *fakeWatcher) AddWatch(path string) error {
	w.watches = append(w.watches, path)
	return nil
}

// Triggers a mock filesystem event.
func (w *fakeWatcher) TriggerEvent(op fsnotify.Op, filename string) {
	w.eventHandler(fsnotify.Event{Op: op, Name: filename})
}
