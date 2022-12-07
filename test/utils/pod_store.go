package utils

import (
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Convenient wrapper around cache.Store that returns list of v1.Pod instead of interface{}.
type PodStore struct {
	cache.Store
	stopCh    chan struct{}
	Reflector *cache.Reflector
}

func NewPodStore(c clientset.Interface, namespace string, label labels.Selector, field fields.Selector) (*PodStore, error) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.LabelSelector = label.String()
			options.FieldSelector = field.String()
			obj, err := c.CoreV1().Pods(namespace).List(options)
			return runtime.Object(obj), err
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.LabelSelector = label.String()
			options.FieldSelector = field.String()
			return c.CoreV1().Pods(namespace).Watch(options)
		},
	}
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	stopCh := make(chan struct{})
	reflector := cache.NewReflector(lw, &v1.Pod{}, store, 0)
	go reflector.Run(stopCh)
	if err := wait.PollImmediate(50*time.Millisecond, 2*time.Minute, func() (bool, error) {
		if len(reflector.LastSyncResourceVersion()) != 0 {
			return true, nil
		}
		return false, nil
	}); err != nil {
		close(stopCh)
		return nil, err
	}
	return &PodStore{Store: store, stopCh: stopCh, Reflector: reflector}, nil
}

func (s *PodStore) List() []*v1.Pod {
	objects := s.Store.List()
	pods := make([]*v1.Pod, 0)
	for _, o := range objects {
		pods = append(pods, o.(*v1.Pod))
	}
	return pods
}

func (s *PodStore) Stop() {
	close(s.stopCh)
}
