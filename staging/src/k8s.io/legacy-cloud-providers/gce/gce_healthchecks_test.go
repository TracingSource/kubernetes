// +build !providerless

package gce

import (
	"testing"

	"k8s.io/api/core/v1"
)

func TestIsAtLeastMinNodesHealthCheckVersion(t *testing.T) {
	testCases := []struct {
		version string
		expect  bool
	}{
		{"v1.7.3", true},
		{"v1.7.2", true},
		{"v1.7.2-alpha.2.597+276d289b90d322", true},
		{"v1.6.0-beta.3.472+831q821c907t31a", false},
		{"v1.5.2", false},
	}

	for _, tc := range testCases {
		if res := isAtLeastMinNodesHealthCheckVersion(tc.version); res != tc.expect {
			t.Errorf("%v: want %v, got %v", tc.version, tc.expect, res)
		}
	}
}

func TestSupportsNodesHealthCheck(t *testing.T) {
	testCases := []struct {
		desc   string
		nodes  []*v1.Node
		expect bool
	}{
		{
			"All nodes support nodes health check",
			[]*v1.Node{
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.7.2",
						},
					},
				},
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.7.2-alpha.2.597+276d289b90d322",
						},
					},
				},
			},
			true,
		},
		{
			"All nodes don't support nodes health check",
			[]*v1.Node{
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.6.0-beta.3.472+831q821c907t31a",
						},
					},
				},
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.5.2",
						},
					},
				},
			},
			false,
		},
		{
			"One node doesn't support nodes health check",
			[]*v1.Node{
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.7.3",
						},
					},
				},
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.7.2-alpha.2.597+276d289b90d322",
						},
					},
				},
				{
					Status: v1.NodeStatus{
						NodeInfo: v1.NodeSystemInfo{
							KubeProxyVersion: "v1.5.2",
						},
					},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		if res := supportsNodesHealthCheck(tc.nodes); res != tc.expect {
			t.Errorf("%v: want %v, got %v", tc.desc, tc.expect, res)
		}
	}
}
