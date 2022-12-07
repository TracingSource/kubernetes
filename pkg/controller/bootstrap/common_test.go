package bootstrap

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core "k8s.io/client-go/testing"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	"k8s.io/kubernetes/pkg/apis/core/helper"
)

func newTokenSecret(tokenID, tokenSecret string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       metav1.NamespaceSystem,
			Name:            bootstrapapi.BootstrapTokenSecretPrefix + tokenID,
			ResourceVersion: "1",
		},
		Type: bootstrapapi.SecretTypeBootstrapToken,
		Data: map[string][]byte{
			bootstrapapi.BootstrapTokenIDKey:     []byte(tokenID),
			bootstrapapi.BootstrapTokenSecretKey: []byte(tokenSecret),
		},
	}
}

func addSecretExpiration(s *v1.Secret, expiration string) {
	s.Data[bootstrapapi.BootstrapTokenExpirationKey] = []byte(expiration)
}

func addSecretSigningUsage(s *v1.Secret, value string) {
	s.Data[bootstrapapi.BootstrapTokenUsageSigningKey] = []byte(value)
}

func verifyActions(t *testing.T, expected, actual []core.Action) {
	for i, a := range actual {
		if len(expected) < i+1 {
			t.Errorf("%d unexpected actions: %s", len(actual)-len(expected), spew.Sdump(actual[i:]))
			break
		}

		e := expected[i]
		if !helper.Semantic.DeepEqual(e, a) {
			t.Errorf("Expected\n\t%s\ngot\n\t%s", spew.Sdump(e), spew.Sdump(a))
			continue
		}
	}

	if len(expected) > len(actual) {
		t.Errorf("%d additional expected actions", len(expected)-len(actual))
		for _, a := range expected[len(actual):] {
			t.Logf("    %s", spew.Sdump(a))
		}
	}
}
