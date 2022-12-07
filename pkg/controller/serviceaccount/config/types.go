package config

// SAControllerConfiguration contains elements describing ServiceAccountController.
type SAControllerConfiguration struct {
	// ServiceAccountKeyFile /etc/kubernetes/pki/sa.key
	//
	// serviceAccountKeyFile is the filename containing a PEM-encoded private RSA key
	// used to sign service account tokens.
	ServiceAccountKeyFile string
	// concurrentSATokenSyncs is the number of service account token syncing operations
	// that will be done concurrently.
	ConcurrentSATokenSyncs int32
	// RootCAFile /etc/kubernetes/pki/ca.crt
	//
	// rootCAFile is the root certificate authority will be included in service
	// account's token secret. This must be a valid PEM-encoded CA bundle.
	RootCAFile string
}
