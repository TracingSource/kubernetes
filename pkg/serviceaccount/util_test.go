package serviceaccount

import (
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsServiceAccountToken(t *testing.T) {

	secretIns := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "token-secret-1",
			Namespace:       "default",
			UID:             "23456",
			ResourceVersion: "1",
			Annotations: map[string]string{
				v1.ServiceAccountNameKey: "default",
				v1.ServiceAccountUIDKey:  "12345",
			},
		},
		Type: v1.SecretTypeServiceAccountToken,
		Data: map[string][]byte{
			"token":     []byte("ABC"),
			"ca.crt":    []byte("CA Data"),
			"namespace": []byte("default"),
		},
	}

	secretTypeMistmatch := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "token-secret-2",
			Namespace:       "default",
			UID:             "23456",
			ResourceVersion: "1",
			Annotations: map[string]string{
				v1.ServiceAccountNameKey: "default",
				v1.ServiceAccountUIDKey:  "12345",
			},
		},
		Type: v1.SecretTypeOpaque,
	}

	saIns := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "default",
			UID:             "12345",
			Namespace:       "default",
			ResourceVersion: "1",
		},
	}

	saInsNameNotEqual := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "non-default",
			UID:             "12345",
			Namespace:       "default",
			ResourceVersion: "1",
		},
	}

	saInsUIDNotEqual := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "default",
			UID:             "67890",
			Namespace:       "default",
			ResourceVersion: "1",
		},
	}

	tests := map[string]struct {
		secret *v1.Secret
		sa     *v1.ServiceAccount
		expect bool
	}{
		"correct service account": {
			secret: secretIns,
			sa:     saIns,
			expect: true,
		},
		"service account name not equal": {
			secret: secretIns,
			sa:     saInsNameNotEqual,
			expect: false,
		},
		"service account uid not equal": {
			secret: secretIns,
			sa:     saInsUIDNotEqual,
			expect: false,
		},
		"service account type not equal": {
			secret: secretTypeMistmatch,
			sa:     saIns,
			expect: false,
		},
	}

	for k, v := range tests {
		actual := IsServiceAccountToken(v.secret, v.sa)
		if actual != v.expect {
			t.Errorf("%s failed, expected %t but received %t", k, v.expect, actual)
		}
	}

}
