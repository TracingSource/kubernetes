package phases

import (
	"testing"

	"k8s.io/component-base/version"
	kubeadmapiv1beta2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta2"
)

func TestSetKubernetesVersion(t *testing.T) {

	ver := version.Get().String()

	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "empty version is processed",
			input:  "",
			output: ver,
		},
		{
			name:   "default version is processed",
			input:  kubeadmapiv1beta2.DefaultKubernetesVersion,
			output: ver,
		},
		{
			name:   "any other version is skipped",
			input:  "v1.12.0",
			output: "v1.12.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &kubeadmapiv1beta2.ClusterConfiguration{KubernetesVersion: test.input}
			SetKubernetesVersion(cfg)
			if cfg.KubernetesVersion != test.output {
				t.Fatalf("expected %q, got %q", test.output, cfg.KubernetesVersion)
			}
		})
	}
}
