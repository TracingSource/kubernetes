package configmapandsecret

import (
	"reflect"
	"testing"

	"k8s.io/cli-runtime/pkg/kustomize/k8sdeps/kv"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
)

func TestKeyValuesFromFileSources(t *testing.T) {
	tests := []struct {
		description string
		sources     []string
		expected    []kv.Pair
	}{
		{
			description: "create kvs from file sources",
			sources:     []string{"files/app-init.ini"},
			expected: []kv.Pair{
				{
					Key:   "app-init.ini",
					Value: "FOO=bar",
				},
			},
		},
	}

	fSys := fs.MakeFakeFS()
	fSys.WriteFile("/files/app-init.ini", []byte("FOO=bar"))
	for _, tc := range tests {
		kvs, err := keyValuesFromFileSources(loader.NewFileLoaderAtRoot(fSys), tc.sources)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(kvs, tc.expected) {
			t.Fatalf("in testcase: %q updated:\n%#v\ndoesn't match expected:\n%#v\n", tc.description, kvs, tc.expected)
		}
	}
}
