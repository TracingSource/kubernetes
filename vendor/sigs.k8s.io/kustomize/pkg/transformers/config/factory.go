package config

import (
	"log"

	"github.com/ghodss/yaml"
	"sigs.k8s.io/kustomize/pkg/ifc"
)

// Factory makes instances of TransformerConfig.
type Factory struct {
	ldr ifc.Loader
}

// MakeTransformerConfig returns a merger of custom config,
// if any, with default config.
func MakeTransformerConfig(
	ldr ifc.Loader, paths []string) (*TransformerConfig, error) {
	t1 := MakeDefaultConfig()
	if len(paths) == 0 {
		return t1, nil
	}
	t2, err := NewFactory(ldr).FromFiles(paths)
	if err != nil {
		return nil, err
	}
	return t1.Merge(t2)
}

func NewFactory(l ifc.Loader) *Factory {
	return &Factory{ldr: l}
}

func (tf *Factory) loader() ifc.Loader {
	if tf.ldr.(ifc.Loader) == nil {
		log.Fatal("no loader")
	}
	return tf.ldr
}

// FromFiles returns a TranformerConfig object from a list of files
func (tf *Factory) FromFiles(
	paths []string) (*TransformerConfig, error) {
	result := &TransformerConfig{}
	for _, path := range paths {
		data, err := tf.loader().Load(path)
		if err != nil {
			return nil, err
		}
		t, err := makeTransformerConfigFromBytes(data)
		if err != nil {
			return nil, err
		}
		result, err = result.Merge(t)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// makeTransformerConfigFromBytes returns a TransformerConfig object from bytes
func makeTransformerConfigFromBytes(data []byte) (*TransformerConfig, error) {
	var t TransformerConfig
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	t.sortFields()
	return &t, nil
}
