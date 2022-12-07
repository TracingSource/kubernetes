package config

import (
	"sync"

	"k8s.io/apimachinery/pkg/util/wait"
)

// Merger 由 pkg/kubelet/config/config.go -> podStorage{} 结构体实现
type Merger interface {
	// Invoked when a change from a source is received. 
	// May also function as an incremental
	// merger if you wish to consume changes incrementally. 
	// Must be reentrant when more than one source is defined.
	Merge(source string, update interface{}) error
}

// MergeFunc implements the Merger interface
type MergeFunc func(source string, update interface{}) error

func (f MergeFunc) Merge(source string, update interface{}) error {
	return f(source, update)
}

// Mux is a class for merging configuration from multiple sources. 
// Changes are pushed via channels and sent to the merge function.
type Mux struct {
	// merger pkg/kubelet/config/config.go -> podStorage{} 结构体
	//
	// Invoked when an update is sent to a source.
	merger Merger

	// Sources and their lock.
	sourceLock sync.RWMutex
	// Maps source names to channels
	sources map[string]chan interface{}
}

// NewMux creates a new mux that can merge changes from multiple sources.
func NewMux(merger Merger) *Mux {
	mux := &Mux{
		sources: make(map[string]chan interface{}),
		merger:  merger,
	}
	return mux
}

// Channel 为目标 source 来源创建新的 channel 通道, 并启动一个协程开始循环处理该通道中的事件.
//
// 	@param source: 所谓 source 来源, 指的是 Pod 的"期望"状态来源, 一般有2种
//         file: 表示 manifests 目录下存放的静态 Pod 信息
//         api: 表示从 apiserver 上同步下来的常规 Pod 信息
//
// caller: 
// 	1. pkg/kubelet/config/config.go -> PodConfig.Channel()
//
// Channel returns a channel where a configuration source
// can send updates of new configurations. Multiple calls with the same
// source will return the same channel. This allows change and state based sources
// to use the same channel. Different source names however will be treated as a
// union.
func (m *Mux) Channel(source string) chan interface{} {
	if len(source) == 0 {
		panic("Channel given an empty name")
	}
	m.sourceLock.Lock()
	defer m.sourceLock.Unlock()
	channel, exists := m.sources[source]
	if exists {
		return channel
	}
	newChannel := make(chan interface{})
	m.sources[source] = newChannel
	go wait.Until(func() { m.listen(source, newChannel) }, 0, wait.NeverStop)
	return newChannel
}

// caller: 
// 	1. Mux.Channel()
func (m *Mux) listen(source string, listenChannel <-chan interface{}) {
	for update := range listenChannel {
		m.merger.Merge(source, update)
	}
}

// Accessor is an interface for retrieving the current merge state.
type Accessor interface {
	// MergedState returns a representation of the current merge state.
	// Must be reentrant when more than one source is defined.
	MergedState() interface{}
}

// AccessorFunc implements the Accessor interface.
type AccessorFunc func() interface{}

func (f AccessorFunc) MergedState() interface{} {
	return f()
}

type Listener interface {
	// OnUpdate is invoked when a change is made to an object.
	OnUpdate(instance interface{})
}

// ListenerFunc receives a representation of the change or object.
type ListenerFunc func(instance interface{})

func (f ListenerFunc) OnUpdate(instance interface{}) {
	f(instance)
}

type Broadcaster struct {
	// Listeners for changes and their lock.
	listenerLock sync.RWMutex
	listeners    []Listener
}

// NewBroadcaster registers a set of listeners that support the Listener interface
// and notifies them all on changes.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

// Add registers listener to receive updates of changes.
func (b *Broadcaster) Add(listener Listener) {
	b.listenerLock.Lock()
	defer b.listenerLock.Unlock()
	b.listeners = append(b.listeners, listener)
}

// Notify notifies all listeners.
func (b *Broadcaster) Notify(instance interface{}) {
	b.listenerLock.RLock()
	listeners := b.listeners
	b.listenerLock.RUnlock()
	for _, listener := range listeners {
		listener.OnUpdate(instance)
	}
}
