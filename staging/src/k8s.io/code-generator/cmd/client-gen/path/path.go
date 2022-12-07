package path

import "strings"

// Vendorless removes the longest match of "*/vendor/" from the front of p.
// It is useful if a package locates in vendor/, e.g.,
// k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/apis/meta/v1, because gengo
// indexes the package with its import path, e.g.,
// k8s.io/apimachinery/pkg/apis/meta/v1,
func Vendorless(p string) string {
	if pos := strings.LastIndex(p, "/vendor/"); pos != -1 {
		return p[pos+len("/vendor/"):]
	}
	return p
}
