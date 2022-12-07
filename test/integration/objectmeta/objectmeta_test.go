package objectmeta

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/kubernetes/test/integration/framework"
)

func TestIgnoreClusterName(t *testing.T) {
	config := framework.NewMasterConfig()
	_, s, closeFn := framework.RunAMaster(config)
	defer closeFn()

	client := clientset.NewForConfigOrDie(&restclient.Config{Host: s.URL, ContentConfig: restclient.ContentConfig{GroupVersion: &schema.GroupVersion{Group: "", Version: "v1"}}})
	ns := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-namespace",
			ClusterName: "cluster-name-to-ignore",
		},
	}
	nsNew, err := client.CoreV1().Namespaces().Create(&ns)
	assert.Nil(t, err)
	assert.Equal(t, ns.Name, nsNew.Name)
	assert.Empty(t, nsNew.ClusterName)

	nsNew, err = client.CoreV1().Namespaces().Update(&ns)
	assert.Nil(t, err)
	assert.Equal(t, ns.Name, nsNew.Name)
	assert.Empty(t, nsNew.ClusterName)
}
