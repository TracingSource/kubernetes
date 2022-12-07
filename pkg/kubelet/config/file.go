package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"k8s.io/klog"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	utilio "k8s.io/utils/io"
)

type podEventType int

const (
	podAdd podEventType = iota
	podModify
	podDelete

	eventBufferLen = 10
)

type watchEvent struct {
	fileName  string
	eventType podEventType
}

type sourceFile struct {
	path           string
	nodeName       types.NodeName
	// period 一般为 20s
	period         time.Duration
	store          cache.Store
	fileKeyMapping map[string]string
	updates        chan<- interface{}
	// watchEvents 对 /etc/kubernetes/manifests 目录的监听事件.
	// 当该目录下的文件发生变动时, 由 run() 方法将其汇总发送到 updates 通道中.
	// 	写入: pkg/kubelet/config/file_linux.go -> sourceFile.produceWatchEvent()
	// 	读取: pkg/kubelet/config/file.go -> sourceFile.run()
	watchEvents    chan *watchEvent
}

// NewSourceFile 开始监听 manifests 目录下的文件, 如果发生变动, 则由 run() 方法将其汇总发送到 updates 通道中.
//
// 	@param path: /etc/kubernetes/manifests 目录
// 	@param nodeName: 当前 kubelet 所在的主机名称
// 	@param updates: 当 manifests 目录下文件发生变动时, 将变动事件发送到此通道中.
//         其值为 pkg/util/config/config.go -> Mux.Channel() 的返回值
//
// caller: 
// 	1. pkg/kubelet/kubelet.go -> makePodSourceConfig() 
//
// NewSourceFile watches a config file for changes.
func NewSourceFile(
	path string, nodeName types.NodeName, period time.Duration, 
	updates chan<- interface{},
) {
	// "github.com/sigma/go-inotify" requires a path without trailing "/"
	path = strings.TrimRight(path, string(os.PathSeparator))

	config := newSourceFile(path, nodeName, period, updates)
	klog.V(1).Infof("Watching path %q", path)
	config.run()
}

// caller: 
// 	1. NewSourceFile()
func newSourceFile(
	path string, nodeName types.NodeName, period time.Duration, 
	updates chan<- interface{},
) *sourceFile {
	send := func(objs []interface{}) {
		// objs 为当前 store 中缓存的 staitcPod 列表.

		var pods []*v1.Pod
		for _, o := range objs {
			pods = append(pods, o.(*v1.Pod))
		}
		// 注意: 这里的 Op 类型为 SET, 在 pkg/kubelet/kubelet.go -> 
		updates <- kubetypes.PodUpdate{
			Pods: pods, Op: kubetypes.SET, Source: kubetypes.FileSource,
		}
	}
	store := cache.NewUndeltaStore(send, cache.MetaNamespaceKeyFunc)
	return &sourceFile{
		path:           path,
		nodeName:       nodeName,
		period:         period,
		store:          store,
		fileKeyMapping: map[string]string{},
		updates:        updates,
		watchEvents:    make(chan *watchEvent, eventBufferLen),
	}
}

// run 开始对 manifests 目录下的文件进行监听. 
// 如有变动, 则立刻将变动内容写入 updates 通道进行上报;
// 如无变动, 则定时将当前目录下的内容写入 store 缓存中;
//
// caller: 
// 	1. NewSourceFile()
func (s *sourceFile) run() {
	listTicker := time.NewTicker(s.period)

	go func() {
		// Read path immediately to speed up startup.
		if err := s.listConfig(); err != nil {
			klog.Errorf("Unable to read config path %q: %v", s.path, err)
		}
		for {
			select {
			case <-listTicker.C:
				if err := s.listConfig(); err != nil {
					klog.Errorf("Unable to read config path %q: %v", s.path, err)
				}
			case e := <-s.watchEvents:
				if err := s.consumeWatchEvent(e); err != nil {
					klog.Errorf("Unable to process watch event: %v", err)
				}
			}
		}
	}()

	s.startWatch()
}

func (s *sourceFile) applyDefaults(pod *api.Pod, source string) error {
	return applyDefaults(pod, source, true, s.nodeName)
}

// caller: 
// 	1. sourceFile.run()
func (s *sourceFile) listConfig() error {
	path := s.path
	statInfo, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// Emit an update with an empty PodList to allow FileSource to be marked as seen
		s.updates <- kubetypes.PodUpdate{
			Pods: []*v1.Pod{}, Op: kubetypes.SET, Source: kubetypes.FileSource,
		}
		return fmt.Errorf("path does not exist, ignoring")
	}

	switch {
	case statInfo.Mode().IsDir():
		pods, err := s.extractFromDir(path)
		if err != nil {
			return err
		}
		if len(pods) == 0 {
			// Emit an update with an empty PodList to allow FileSource to be marked as seen
			s.updates <- kubetypes.PodUpdate{
				Pods: pods, Op: kubetypes.SET, Source: kubetypes.FileSource,
			}
			return nil
		}
		return s.replaceStore(pods...)

	case statInfo.Mode().IsRegular():
		pod, err := s.extractFromFile(path)
		if err != nil {
			return err
		}
		return s.replaceStore(pod)

	default:
		return fmt.Errorf("path is not a directory or file")
	}
}

// Get as many pod manifests as we can from a directory. Return an error if and only if something
// prevented us from reading anything at all. Do not return an error if only some files
// were problematic.
func (s *sourceFile) extractFromDir(name string) ([]*v1.Pod, error) {
	dirents, err := filepath.Glob(filepath.Join(name, "[^.]*"))
	if err != nil {
		return nil, fmt.Errorf("glob failed: %v", err)
	}

	pods := make([]*v1.Pod, 0)
	if len(dirents) == 0 {
		return pods, nil
	}

	sort.Strings(dirents)
	for _, path := range dirents {
		statInfo, err := os.Stat(path)
		if err != nil {
			klog.Errorf("Can't get metadata for %q: %v", path, err)
			continue
		}

		switch {
		case statInfo.Mode().IsDir():
			klog.Errorf("Not recursing into manifest path %q", path)
		case statInfo.Mode().IsRegular():
			pod, err := s.extractFromFile(path)
			if err != nil {
				if !os.IsNotExist(err) {
					klog.Errorf("Can't process manifest file %q: %v", path, err)
				}
			} else {
				pods = append(pods, pod)
			}
		default:
			klog.Errorf("Manifest path %q is not a directory or file: %v", path, statInfo.Mode())
		}
	}
	return pods, nil
}

// extractFromFile 读取 manifests 目录下的指定文件, 并将其转换成 Pod 对象.
//
// extractFromFile parses a file for Pod configuration information.
func (s *sourceFile) extractFromFile(filename string) (pod *v1.Pod, err error) {
	klog.V(3).Infof("Reading config file %q", filename)
	defer func() {
		if err == nil && pod != nil {
			objKey, keyErr := cache.MetaNamespaceKeyFunc(pod)
			if keyErr != nil {
				err = keyErr
				return
			}
			s.fileKeyMapping[filename] = objKey
		}
	}()

	file, err := os.Open(filename)
	if err != nil {
		return pod, err
	}
	defer file.Close()

	data, err := utilio.ReadAtMost(file, maxConfigLength)
	if err != nil {
		return pod, err
	}

	defaultFn := func(pod *api.Pod) error {
		return s.applyDefaults(pod, filename)
	}

	parsed, pod, podErr := tryDecodeSinglePod(data, defaultFn)
	if parsed {
		if podErr != nil {
			return pod, podErr
		}
		return pod, nil
	}

	return pod, fmt.Errorf("%v: couldn't parse as pod(%v), please check config file", filename, podErr)
}

func (s *sourceFile) replaceStore(pods ...*v1.Pod) (err error) {
	objs := []interface{}{}
	for _, pod := range pods {
		objs = append(objs, pod)
	}
	return s.store.Replace(objs, "")
}
