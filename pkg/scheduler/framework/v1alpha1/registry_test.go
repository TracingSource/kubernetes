package v1alpha1

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/apis/config"
	"k8s.io/kubernetes/pkg/scheduler/apis/config/scheme"
)

func TestDecodeInto(t *testing.T) {
	type PluginFooConfig struct {
		FooTest string `json:"foo_test,omitempty"`
	}
	tests := []struct {
		name            string
		schedulerConfig string
		expeted         PluginFooConfig
	}{
		{
			name: "test decode for JSON config",
			schedulerConfig: `{
				"kind": "KubeSchedulerConfiguration",
				"apiVersion": "kubescheduler.config.k8s.io/v1alpha1",
				"plugins": {
				"permit": {
						"enabled": [
							{
								"name": "foo"
							}
						]
					}
				},
				"pluginConfig": [
					{
						"name": "foo",
						"args": {
							"foo_test": "test decode"
						}
					}
				]
			}`,
			expeted: PluginFooConfig{
				FooTest: "test decode",
			},
		},
		{
			name: "test decode for YAML config",
			schedulerConfig: `
apiVersion: kubescheduler.config.k8s.io/v1alpha1
kind: KubeSchedulerConfiguration
plugins:
  permit:
    enabled:
      - name: foo
pluginConfig:
  - name: foo
    args:
      foo_test: "test decode"`,
			expeted: PluginFooConfig{
				FooTest: "test decode",
			},
		},
	}
	for i, test := range tests {
		schedulerConf, err := loadConfig([]byte(test.schedulerConfig))
		if err != nil {
			t.Errorf("Test #%v(%s): failed to load scheduler config: %v", i, test.name, err)
		}
		var pluginFooConf PluginFooConfig
		if err := DecodeInto(&schedulerConf.PluginConfig[0].Args, &pluginFooConf); err != nil {
			t.Errorf("Test #%v(%s): failed to decode args %+v: %v",
				i, test.name, schedulerConf.PluginConfig[0].Args, err)
		}
		if !reflect.DeepEqual(pluginFooConf, test.expeted) {
			t.Errorf("Test #%v(%s): failed to decode plugin config, expected: %+v, got: %+v",
				i, test.name, test.expeted, pluginFooConf)
		}
	}
}

func loadConfig(data []byte) (*config.KubeSchedulerConfiguration, error) {
	configObj := &config.KubeSchedulerConfiguration{}
	if err := runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), data, configObj); err != nil {
		return nil, err
	}

	return configObj, nil
}
