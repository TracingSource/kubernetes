package kubeadm

import (
	yaml "gopkg.in/yaml.v2"
	authv1 "k8s.io/api/authorization/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	kubeadmConfigName                             = "kubeadm-config"
	kubeadmConfigRoleName                         = "kubeadm:nodes-kubeadm-config"
	kubeadmConfigRoleBindingName                  = kubeadmConfigRoleName
	kubeadmConfigClusterConfigurationConfigMapKey = "ClusterConfiguration"
	kubeadmConfigClusterStatusConfigMapKey        = "ClusterStatus"
)

var (
	kubeadmConfigConfigMapResource = &authv1.ResourceAttributes{
		Namespace: kubeSystemNamespace,
		Name:      kubeadmConfigName,
		Resource:  "configmaps",
		Verb:      "get",
	}
)

// Define container for all the test specification aimed at verifying
// that kubeadm creates the cluster-info ConfigMap, that it is properly configured
// and that all the related RBAC rules are in place
var _ = Describe("kubeadm-config ConfigMap", func() {

	// Get an instance of the k8s test framework
	f := framework.NewDefaultFramework("kubeadm-config")

	// Tests in this container are not expected to create new objects in the cluster
	// so we are disabling the creation of a namespace in order to get a faster execution
	f.SkipNamespaceCreation = true

	ginkgo.It("should exist and be properly configured", func() {
		cm := GetConfigMap(f.ClientSet, kubeSystemNamespace, kubeadmConfigName)

		gomega.Expect(cm.Data).To(gomega.HaveKey(kubeadmConfigClusterConfigurationConfigMapKey))
		gomega.Expect(cm.Data).To(gomega.HaveKey(kubeadmConfigClusterStatusConfigMapKey))

		m := unmarshalYaml(cm.Data[kubeadmConfigClusterStatusConfigMapKey])
		if _, ok := m["apiEndpoints"]; ok {
			d := m["apiEndpoints"].(map[interface{}]interface{})
			// get all control-plane nodes
			controlPlanes := getControlPlaneNodes(f.ClientSet)

			// checks that all the control-plane nodes are in the apiEndpoints list
			for _, cp := range controlPlanes.Items {
				if _, ok := d[cp.Name]; !ok {
					framework.Failf("failed to get apiEndpoints for control-plane %s in %s", cp.Name, kubeadmConfigClusterStatusConfigMapKey)
				}
			}
		} else {
			framework.Failf("failed to get apiEndpoints from %s", kubeadmConfigClusterStatusConfigMapKey)
		}
	})

	ginkgo.It("should have related Role and RoleBinding", func() {
		ExpectRole(f.ClientSet, kubeSystemNamespace, kubeadmConfigRoleName)
		ExpectRoleBinding(f.ClientSet, kubeSystemNamespace, kubeadmConfigRoleBindingName)
	})

	ginkgo.It("should be accessible for bootstrap tokens", func() {
		ExpectSubjectHasAccessToResource(f.ClientSet,
			rbacv1.GroupKind, bootstrapTokensGroup,
			kubeadmConfigConfigMapResource,
		)
	})

	ginkgo.It("should be accessible for nodes", func() {
		ExpectSubjectHasAccessToResource(f.ClientSet,
			rbacv1.GroupKind, nodesGroup,
			kubeadmConfigConfigMapResource,
		)
	})
})

func getClusterConfiguration(c clientset.Interface) map[interface{}]interface{} {
	cm := GetConfigMap(c, kubeSystemNamespace, kubeadmConfigName)

	gomega.Expect(cm.Data).To(gomega.HaveKey(kubeadmConfigClusterConfigurationConfigMapKey))

	return unmarshalYaml(cm.Data[kubeadmConfigClusterConfigurationConfigMapKey])
}

func unmarshalYaml(data string) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		framework.Failf("error parsing %s ConfigMap: %v", kubeadmConfigName, err)
	}
	return m
}
