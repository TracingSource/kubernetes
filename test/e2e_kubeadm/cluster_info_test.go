package kubeadm

import (
	authv1 "k8s.io/api/authorization/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	clusterInfoConfigMapName   = "cluster-info"
	clusterInfoRoleName        = "kubeadm:bootstrap-signer-clusterinfo"
	clusterInfoRoleBindingName = clusterInfoRoleName
)

var (
	clusterInfoConfigMapResource = &authv1.ResourceAttributes{
		Namespace: kubePublicNamespace,
		Name:      clusterInfoConfigMapName,
		Resource:  "configmaps",
		Verb:      "get",
	}
)

// Define container for all the test specification aimed at verifying
// that kubeadm creates the cluster-info ConfigMap, that it is properly configured
// and that all the related RBAC rules are in place
var _ = Describe("cluster-info ConfigMap", func() {

	// Get an instance of the k8s test framework
	f := framework.NewDefaultFramework("cluster-info")

	// Tests in this container are not expected to create new objects in the cluster
	// so we are disabling the creation of a namespace in order to get a faster execution
	f.SkipNamespaceCreation = true

	ginkgo.It("should exist and be properly configured", func() {
		// Nb. this is technically implemented a part of the bootstrap-token phase
		cm := GetConfigMap(f.ClientSet, kubePublicNamespace, clusterInfoConfigMapName)

		gomega.Expect(cm.Data).To(gomega.HaveKey(gomega.HavePrefix(bootstrapapi.JWSSignatureKeyPrefix)))
		gomega.Expect(cm.Data).To(gomega.HaveKey(bootstrapapi.KubeConfigKey))

		//TODO: What else? server?
	})

	ginkgo.It("should have related Role and RoleBinding", func() {
		// Nb. this is technically implemented a part of the bootstrap-token phase
		ExpectRole(f.ClientSet, kubePublicNamespace, clusterInfoRoleName)
		ExpectRoleBinding(f.ClientSet, kubePublicNamespace, clusterInfoRoleBindingName)
	})

	ginkgo.It("should be accessible for anonymous", func() {
		ExpectSubjectHasAccessToResource(f.ClientSet,
			rbacv1.UserKind, anonymousUser,
			clusterInfoConfigMapResource,
		)
	})
})
