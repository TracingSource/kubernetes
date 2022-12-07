package kubelet

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
)

func TestCreateConfigMap(t *testing.T) {
	nodeName := "fake-node"
	client := fake.NewSimpleClientset()
	k8sVersionStr := constants.CurrentKubernetesVersion.String()
	cfg := &kubeletconfigv1beta1.KubeletConfiguration{}

	client.PrependReactor("get", "nodes", func(action core.Action) (bool, runtime.Object, error) {
		return true, &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: nodeName,
			},
			Spec: v1.NodeSpec{},
		}, nil
	})
	client.PrependReactor("create", "roles", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "rolebindings", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})

	if err := CreateConfigMap(cfg, k8sVersionStr, client); err != nil {
		t.Errorf("CreateConfigMap: unexpected error %v", err)
	}
}

func TestCreateConfigMapRBACRules(t *testing.T) {
	client := fake.NewSimpleClientset()
	client.PrependReactor("create", "roles", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})
	client.PrependReactor("create", "rolebindings", func(action core.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})

	if err := createConfigMapRBACRules(client, version.MustParseSemantic("v1.11.0")); err != nil {
		t.Errorf("createConfigMapRBACRules: unexpected error %v", err)
	}
}
