package kubeconfig

/*

	PHASE: KUBECONFIG

	INPUTS:
		From InitConfiguration
			The API Server endpoint (AdvertiseAddress + BindPort) is required so the kubeconfig file knows where to find the control-plane
			The KubernetesDir path is required for knowing where to put the kubeconfig files
			The PKIPath is required for knowing where all certificates should be stored

	OUTPUTS:
		Files to KubernetesDir (default /etc/kubernetes):
		 - admin.conf
		 - kubelet.conf
		 - scheduler.conf
		 - controller-manager.conf
*/
