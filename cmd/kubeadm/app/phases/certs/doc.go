package certs

/*

	PHASE: CERTIFICATES

	INPUTS:
		From InitConfiguration
			.API.AdvertiseAddress is an optional parameter that can be passed for an extra addition to the SAN IPs
			.APIServer.CertSANs is an optional parameter for adding DNS names and IPs to the API Server serving cert SAN
			.Etcd.Local.ServerCertSANs is an optional parameter for adding DNS names and IPs to the etcd serving cert SAN
			.Etcd.Local.PeerCertSANs is an optional parameter for adding DNS names and IPs to the etcd peer cert SAN
			.Networking.DNSDomain is needed for knowing which DNS name the internal Kubernetes service has
			.Networking.ServiceSubnet is needed for knowing which IP the internal Kubernetes service is going to point to
			.CertificatesDir is required for knowing where all certificates should be stored

	OUTPUTS:
		Files to .CertificatesDir (default /etc/kubernetes/pki):
		 - ca.crt
		 - ca.key
		 - apiserver.crt
		 - apiserver.key
		 - apiserver-kubelet-client.crt
		 - apiserver-kubelet-client.key
		 - apiserver-etcd-client.crt
		 - apiserver-etcd-client.key
		 - etcd/ca.crt
		 - etcd/ca.key
		 - etcd/server.crt
		 - etcd/server.key
		 - etcd/peer.crt
		 - etcd/peer.key
		 - etcd/healthcheck-client.crt
		 - etcd/healthcheck-client.key
		 - sa.pub
		 - sa.key
		 - front-proxy-ca.crt
		 - front-proxy-ca.key
		 - front-proxy-client.crt
		 - front-proxy-client.key

*/
