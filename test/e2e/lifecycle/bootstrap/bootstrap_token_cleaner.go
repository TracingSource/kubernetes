package bootstrap

import (
	"time"

	"github.com/onsi/ginkgo"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/lifecycle"
)

var secretNeedClean string
var _ = lifecycle.SIGDescribe("[Feature:BootstrapTokens]", func() {

	var c clientset.Interface

	f := framework.NewDefaultFramework("bootstrap-token-cleaner")

	ginkgo.BeforeEach(func() {
		c = f.ClientSet
	})

	ginkgo.AfterEach(func() {
		if len(secretNeedClean) > 0 {
			ginkgo.By("delete the bootstrap token secret")
			err := c.CoreV1().Secrets(metav1.NamespaceSystem).Delete(secretNeedClean, &metav1.DeleteOptions{})
			secretNeedClean = ""
			framework.ExpectNoError(err)
		}
	})
	ginkgo.It("should delete the token secret when the secret expired", func() {
		ginkgo.By("create a new expired bootstrap token secret")
		tokenId, err := GenerateTokenId()
		framework.ExpectNoError(err)
		tokenSecret, err := GenerateTokenSecret()
		framework.ExpectNoError(err)

		secret := newTokenSecret(tokenId, tokenSecret)
		addSecretExpiration(secret, TimeStringFromNow(-time.Hour))
		_, err = c.CoreV1().Secrets(metav1.NamespaceSystem).Create(secret)

		framework.ExpectNoError(err)

		ginkgo.By("wait for the bootstrap token secret be deleted")
		err = WaitForBootstrapTokenSecretToDisappear(c, tokenId)
		framework.ExpectNoError(err)
	})

	ginkgo.It("should not delete the token secret when the secret is not expired", func() {
		ginkgo.By("create a new expired bootstrap token secret")
		tokenId, err := GenerateTokenId()
		framework.ExpectNoError(err)
		tokenSecret, err := GenerateTokenSecret()
		framework.ExpectNoError(err)
		secret := newTokenSecret(tokenId, tokenSecret)
		addSecretExpiration(secret, TimeStringFromNow(time.Hour))
		_, err = c.CoreV1().Secrets(metav1.NamespaceSystem).Create(secret)
		secretNeedClean = bootstrapapi.BootstrapTokenSecretPrefix + tokenId
		framework.ExpectNoError(err)

		ginkgo.By("wait for the bootstrap token secret not be deleted")
		err = WaitForBootstrapTokenSecretNotDisappear(c, tokenId, 20*time.Second)
		framework.ExpectNoError(err)
	})
})
