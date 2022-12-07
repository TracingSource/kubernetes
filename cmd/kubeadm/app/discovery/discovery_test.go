package discovery

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

func TestFor(t *testing.T) {
	tests := []struct {
		name   string
		d      kubeadm.JoinConfiguration
		expect bool
	}{
		{
			name:   "default Discovery",
			d:      kubeadm.JoinConfiguration{},
			expect: false,
		},
		{
			name: "file Discovery with a path",
			d: kubeadm.JoinConfiguration{
				Discovery: kubeadm.Discovery{
					File: &kubeadm.FileDiscovery{
						KubeConfigPath: "notnil",
					},
				},
			},
			expect: false,
		},
		{
			name: "file Discovery with an url",
			d: kubeadm.JoinConfiguration{
				Discovery: kubeadm.Discovery{
					File: &kubeadm.FileDiscovery{
						KubeConfigPath: "https://localhost",
					},
				},
			},
			expect: false,
		},
		{
			name: "BootstrapTokenDiscovery",
			d: kubeadm.JoinConfiguration{
				Discovery: kubeadm.Discovery{
					BootstrapToken: &kubeadm.BootstrapTokenDiscovery{
						Token: "foo.bar@foobar",
					},
				},
			},
			expect: false,
		},
	}
	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			config := rt.d
			config.Discovery.Timeout = &metav1.Duration{Duration: 5 * time.Minute}
			_, actual := For(&config)
			if (actual == nil) != rt.expect {
				t.Errorf(
					"failed For:\n\texpected: %t\n\t  actual: %t",
					rt.expect,
					(actual == nil),
				)
			}
		})
	}
}
