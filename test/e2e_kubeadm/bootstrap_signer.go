package kubeadm

import (
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
)

const (
	bootstrapTokensSignerRoleName = "system:controller:bootstrap-signer"
)

// Define container for all the test specification aimed at verifying
// that kubeadm creates the bootstrap signer
var _ = Describe("bootstrap signer", func() {

	// Get an instance of the k8s test framework
	f := framework.NewDefaultFramework("bootstrap token")

	// Tests in this container are not expected to create new objects in the cluster
	// so we are disabling the creation of a namespace in order to get a faster execution
	f.SkipNamespaceCreation = true

	ginkgo.It("should be active", func() {
		//NB. this is technically implemented a part of the control-plane phase
		//    and more specifically if the controller manager is properly configured,
		//    the bootstrapsigner controller is activated and the system:controller:bootstrap-signer
		//    group will be automatically created
		ExpectRole(f.ClientSet, kubeSystemNamespace, bootstrapTokensSignerRoleName)
	})
})
