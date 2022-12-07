package kubelet

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
)

func TestEnableDynamicConfigForNode(t *testing.T) {
	nodeName := "fake-node"
	client := fake.NewSimpleClientset()
	client.PrependReactor("get", "nodes", func(action core.Action) (bool, runtime.Object, error) {
		return true, &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   nodeName,
				Labels: map[string]string{v1.LabelHostname: nodeName},
			},
			Spec: v1.NodeSpec{
				ConfigSource: &v1.NodeConfigSource{
					ConfigMap: &v1.ConfigMapNodeConfigSource{
						UID: "",
					},
				},
			},
		}, nil
	})
	client.PrependReactor("get", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
		return true, &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kubelet-config-1.11",
				Namespace: metav1.NamespaceSystem,
				UID:       "fake-uid",
			},
		}, nil
	})
	client.PrependReactor("patch", "nodes", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})

	if err := EnableDynamicConfigForNode(client, nodeName, version.MustParseSemantic("v1.11.0")); err != nil {
		t.Errorf("UpdateNodeWithConfigMap: unexpected error %v", err)
	}
}
