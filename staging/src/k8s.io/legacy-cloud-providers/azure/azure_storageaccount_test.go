// +build !providerless

package azure

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
)

func TestGetStorageAccessKeys(t *testing.T) {
	cloud := &Cloud{}
	fake := newFakeStorageAccountClient()
	cloud.StorageAccountClient = fake
	value := "foo bar"

	tests := []struct {
		results     storage.AccountListKeysResult
		expectedKey string
		expectErr   bool
		err         error
	}{
		{storage.AccountListKeysResult{}, "", true, nil},
		{
			storage.AccountListKeysResult{
				Keys: &[]storage.AccountKey{
					{Value: &value},
				},
			},
			"bar",
			false,
			nil,
		},
		{
			storage.AccountListKeysResult{
				Keys: &[]storage.AccountKey{
					{},
					{Value: &value},
				},
			},
			"bar",
			false,
			nil,
		},
		{storage.AccountListKeysResult{}, "", true, fmt.Errorf("test error")},
	}

	for _, test := range tests {
		expectedKey := test.expectedKey
		fake.Keys = test.results
		fake.Err = test.err
		key, err := cloud.GetStorageAccesskey("acct", "rg")
		if test.expectErr && err == nil {
			t.Errorf("Unexpected non-error")
			continue
		}
		if !test.expectErr && err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		if key != expectedKey {
			t.Errorf("expected: %s, saw %s", expectedKey, key)
		}
	}
}
