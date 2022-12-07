package nodetaint

import (
	"context"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/admission"
	api "k8s.io/kubernetes/pkg/apis/core"
)

const (
	// PluginName is the name of the plugin.
	PluginName = "TaintNodesByCondition"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewPlugin(), nil
	})
}

// NewPlugin creates a new NodeTaint admission plugin.
// This plugin identifies requests from nodes
func NewPlugin() *Plugin {
	return &Plugin{
		Handler: admission.NewHandler(admission.Create),
	}
}

// Plugin holds state for and implements the admission plugin.
type Plugin struct {
	*admission.Handler
}

var (
	_ = admission.Interface(&Plugin{})
)

var (
	nodeResource = api.Resource("nodes")
)

// Admit is the main function that checks node identity and adds taints as needed.
func (p *Plugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	// Our job is just to taint nodes.
	if a.GetResource().GroupResource() != nodeResource || a.GetSubresource() != "" {
		return nil
	}

	node, ok := a.GetObject().(*api.Node)
	if !ok {
		return admission.NewForbidden(a, fmt.Errorf("unexpected type %T", a.GetObject()))
	}

	// Taint node with NotReady taint at creation. This is needed to make sure
	// that nodes are added to the cluster with the NotReady taint. Otherwise,
	// a new node may receive the taint with some delay causing pods to be
	// scheduled on a not-ready node. Node controller will remove the taint
	// when the node becomes ready.
	addNotReadyTaint(node)
	return nil
}

func addNotReadyTaint(node *api.Node) {
	notReadyTaint := api.Taint{
		Key:    v1.TaintNodeNotReady,
		Effect: api.TaintEffectNoSchedule,
	}
	for _, taint := range node.Spec.Taints {
		if taint.MatchTaint(notReadyTaint) {
			// the taint already exists.
			return
		}
	}
	node.Spec.Taints = append(node.Spec.Taints, notReadyTaint)
}
