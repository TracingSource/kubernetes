// Package config provides decoupling between various configuration sources (etcd, files,...) and
// the pieces that actually care about them (loadbalancer, proxy). Config takes 1 or more
// configuration sources and allows for incremental (add/remove) and full replace (set)
// changes from each of the sources, then creates a union of the configuration and provides
// a unified view for both service handlers as well as endpoint handlers. There is no attempt
// to resolve conflicts of any sort. Basic idea is that each configuration source gets a channel
// from the Config service and pushes updates to it via that channel. Config then keeps track of
// incremental & replace changes and distributes them to listeners as appropriate.
package config // import "k8s.io/kubernetes/pkg/proxy/config"
