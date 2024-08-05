package config

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

// NodeHandler is an abstract interface of objects which receive
// notifications about node object changes.
type NodeHandler interface {
	// OnNodeAdd is called whenever creation of new node object
	// is observed.
	OnNodeAdd(node *v1.Node)
	// OnNodeUpdate is called whenever modification of an existing
	// node object is observed.
	OnNodeUpdate(oldNode, node *v1.Node)
	// OnNodeDelete is called whever deletion of an existing node
	// object is observed.
	OnNodeDelete(node *v1.Node)
	// OnNodeSynced is called once all the initial event handlers were
	// called and the state is fully propagated to local cache.
	OnNodeSynced()
}

// NoopNodeHandler is a noop handler for proxiers that have not yet
// implemented a full NodeHandler.
type NoopNodeHandler struct{}

// OnNodeAdd is a noop handler for Node creates.
func (*NoopNodeHandler) OnNodeAdd(node *v1.Node) {}

// OnNodeUpdate is a noop handler for Node updates.
func (*NoopNodeHandler) OnNodeUpdate(oldNode, node *v1.Node) {}

// OnNodeDelete is a noop handler for Node deletes.
func (*NoopNodeHandler) OnNodeDelete(node *v1.Node) {}

// OnNodeSynced is a noop handler for Node syncs.
func (*NoopNodeHandler) OnNodeSynced() {}

// NodeConfig tracks a set of node configurations.
// It accepts "set", "add" and "remove" operations of node via channels, and invokes registered handlers on change.
type NodeConfig struct {
	listerSynced  cache.InformerSynced
	eventHandlers []NodeHandler
}

// NewNodeConfig creates a new NodeConfig.
func NewNodeConfig(nodeInformer coreinformers.NodeInformer, resyncPeriod time.Duration) *NodeConfig {
	result := &NodeConfig{
		listerSynced: nodeInformer.Informer().HasSynced,
	}

	nodeInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    result.handleAddNode,
			UpdateFunc: result.handleUpdateNode,
			DeleteFunc: result.handleDeleteNode,
		},
		resyncPeriod,
	)

	return result
}

// RegisterEventHandler registers a handler which is called on every node change.
func (c *NodeConfig) RegisterEventHandler(handler NodeHandler) {
	c.eventHandlers = append(c.eventHandlers, handler)
}

// Run starts the goroutine responsible for calling registered handlers.
func (c *NodeConfig) Run(stopCh <-chan struct{}) {
	klog.Info("Starting node config controller")

	if !cache.WaitForNamedCacheSync("node config", stopCh, c.listerSynced) {
		return
	}

	for i := range c.eventHandlers {
		klog.V(3).Infof("Calling handler.OnNodeSynced()")
		c.eventHandlers[i].OnNodeSynced()
	}
}

func (c *NodeConfig) handleAddNode(obj interface{}) {
	node, ok := obj.(*v1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
		return
	}
	for i := range c.eventHandlers {
		klog.V(4).Infof("Calling handler.OnNodeAdd")
		c.eventHandlers[i].OnNodeAdd(node)
	}
}

func (c *NodeConfig) handleUpdateNode(oldObj, newObj interface{}) {
	oldNode, ok := oldObj.(*v1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", oldObj))
		return
	}
	node, ok := newObj.(*v1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", newObj))
		return
	}
	for i := range c.eventHandlers {
		klog.V(5).Infof("Calling handler.OnNodeUpdate")
		c.eventHandlers[i].OnNodeUpdate(oldNode, node)
	}
}

func (c *NodeConfig) handleDeleteNode(obj interface{}) {
	node, ok := obj.(*v1.Node)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
		if node, ok = tombstone.Obj.(*v1.Node); !ok {
			utilruntime.HandleError(fmt.Errorf("unexpected object type: %v", obj))
			return
		}
	}
	for i := range c.eventHandlers {
		klog.V(4).Infof("Calling handler.OnNodeDelete")
		c.eventHandlers[i].OnNodeDelete(node)
	}
}
