// Package config implements the pod configuration readers.
package config

import (
	"sync"

	"k8s.io/apimachinery/pkg/util/sets"
)

// SourcesReadyFn is function that returns true if the specified sources have been seen.
type SourcesReadyFn func(sourcesSeen sets.String) bool

// SourcesReady tracks the set of configured sources seen by the kubelet.
type SourcesReady interface {
	// AddSource adds the specified source to the set of sources managed.
	AddSource(source string)
	// AllReady returns true if the currently configured sources have all been seen.
	AllReady() bool
}

// NewSourcesReady returns a SourcesReady with the specified function.
func NewSourcesReady(sourcesReadyFn SourcesReadyFn) SourcesReady {
	return &sourcesImpl{
		sourcesSeen:    sets.NewString(),
		sourcesReadyFn: sourcesReadyFn,
	}
}

// sourcesImpl implements SourcesReady.  It is thread-safe.
type sourcesImpl struct {
	// lock protects access to sources seen.
	lock sync.RWMutex
	// set of sources seen.
	sourcesSeen sets.String
	// sourcesReady is a function that evaluates if the sources are ready.
	sourcesReadyFn SourcesReadyFn
}

// Add adds the specified source to the set of sources managed.
func (s *sourcesImpl) AddSource(source string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sourcesSeen.Insert(source)
}

// AllReady returns true if each configured source is ready.
func (s *sourcesImpl) AllReady() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.sourcesReadyFn(s.sourcesSeen)
}
