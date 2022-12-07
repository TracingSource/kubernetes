// Package kubeapiserver holds code that is common to both the kube-apiserver
// and the federation-apiserver, but isn't part of a generic API server.
// For instance, the non-delegated authorization options are used by those two
// servers, but no generic API server is likely to use them.
package kubeapiserver
