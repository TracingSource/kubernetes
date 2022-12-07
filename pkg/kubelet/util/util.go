package util

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FromApiserverCache modifies <opts> so that the GET request will
// be served from apiserver cache instead of from etcd.
func FromApiserverCache(opts *metav1.GetOptions) {
	opts.ResourceVersion = "0"
}
