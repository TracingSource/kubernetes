// +build !linux

package kubelet

// Do nothing.
func (kl *Kubelet) initNetworkUtil() {}
