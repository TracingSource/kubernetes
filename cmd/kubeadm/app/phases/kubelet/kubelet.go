package kubelet

import (
	"fmt"

	"k8s.io/kubernetes/cmd/kubeadm/app/util/initsystem"
)

// TryStartKubelet attempts to bring up kubelet service
func TryStartKubelet() {
	// If we notice that the kubelet service is inactive, try to start it
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		fmt.Println("[kubelet-start] no supported init system detected, won't make sure the kubelet is running properly.")
		return
	}

	if !initSystem.ServiceExists("kubelet") {
		fmt.Println("[kubelet-start] couldn't detect a kubelet service, can't make sure the kubelet is running properly.")
	}

	// This runs "systemctl daemon-reload && systemctl restart kubelet"
	if err := initSystem.ServiceRestart("kubelet"); err != nil {
		fmt.Printf("[kubelet-start] WARNING: unable to start the kubelet service: [%v]\n", err)
		fmt.Printf("[kubelet-start] Please ensure kubelet is reloaded and running manually.\n")
	}
}

// TryStopKubelet attempts to bring down the kubelet service momentarily
func TryStopKubelet() {
	// If we notice that the kubelet service is inactive, try to start it
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		fmt.Println("[kubelet-start] no supported init system detected, won't make sure the kubelet not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists("kubelet") {
		fmt.Println("[kubelet-start] couldn't detect a kubelet service, can't make sure the kubelet not running for a short period of time while setting up configuration for it.")
	}

	// This runs "systemctl daemon-reload && systemctl stop kubelet"
	if err := initSystem.ServiceStop("kubelet"); err != nil {
		fmt.Printf("[kubelet-start] WARNING: unable to stop the kubelet service momentarily: [%v]\n", err)
	}
}

// TryRestartKubelet attempts to restart the kubelet service
func TryRestartKubelet() {
	// If we notice that the kubelet service is inactive, try to start it
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		fmt.Println("[kubelet-start] no supported init system detected, won't make sure the kubelet not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists("kubelet") {
		fmt.Println("[kubelet-start] couldn't detect a kubelet service, can't make sure the kubelet not running for a short period of time while setting up configuration for it.")
	}

	// This runs "systemctl daemon-reload && systemctl stop kubelet"
	if err := initSystem.ServiceRestart("kubelet"); err != nil {
		fmt.Printf("[kubelet-start] WARNING: unable to restart the kubelet service momentarily: [%v]\n", err)
	}
}
