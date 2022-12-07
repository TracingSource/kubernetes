package node

import (
	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
)

// TODO(mattmoyer): Move CreateNewTokens, UpdateOrCreateTokens out of this package
// to client-go for a generic abstraction and client for a Bootstrap Token

// CreateNewTokens tries to create a token and fails if one with the same ID already exists
func CreateNewTokens(client clientset.Interface, tokens []kubeadmapi.BootstrapToken) error {
	return UpdateOrCreateTokens(client, true, tokens)
}

// caller:
// 	1. cmd/kubeadm/app/cmd/phases/init/bootstraptoken.go -> runBootstrapToken()
// 	在 kubeadm init 过程中被调用.
//
// UpdateOrCreateTokens attempts to update a token with the given ID,
// or create if it does not already exist.
func UpdateOrCreateTokens(
	client clientset.Interface, failIfExists bool, tokens []kubeadmapi.BootstrapToken,
) error {
	for _, token := range tokens {
		secretName := bootstraputil.BootstrapTokenSecretName(token.Token.ID)
		secret, err := client.CoreV1().Secrets(metav1.NamespaceSystem).Get(secretName, metav1.GetOptions{})
		if secret != nil && err == nil && failIfExists {
			return errors.Errorf("a token with id %q already exists", token.Token.ID)
		}

		updatedOrNewSecret := token.ToSecret()
		// Try to create or update the token with an exponential backoff
		err = apiclient.TryRunCommand(func() error {
			err := apiclient.CreateOrUpdateSecret(client, updatedOrNewSecret)
			if err != nil {
				return errors.Wrapf(err, "failed to create or update bootstrap token with name %s", secretName)
			}
			return nil
		}, 5)
		if err != nil {
			return err
		}
	}
	return nil
}
