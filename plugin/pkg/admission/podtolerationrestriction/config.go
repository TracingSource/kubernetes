package podtolerationrestriction

import (
	"fmt"
	"io"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	internalapi "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction"
	"k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/install"
	versionedapi "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/v1alpha1"
	"k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction/apis/podtolerationrestriction/validation"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func init() {
	install.Install(scheme)
}

// LoadConfiguration loads the provided configuration.
func loadConfiguration(config io.Reader) (*internalapi.Configuration, error) {
	// if no config is provided, return a default configuration
	if config == nil {
		externalConfig := &versionedapi.Configuration{}
		scheme.Default(externalConfig)
		internalConfig := &internalapi.Configuration{}
		if err := scheme.Convert(externalConfig, internalConfig, nil); err != nil {
			return nil, err
		}
		return internalConfig, nil
	}
	// we have a config so parse it.
	data, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}
	decoder := codecs.UniversalDecoder()
	decodedObj, err := runtime.Decode(decoder, data)
	if err != nil {
		return nil, err
	}
	externalConfig, ok := decodedObj.(*internalapi.Configuration)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", decodedObj)
	}

	if err := validation.ValidateConfiguration(externalConfig); err != nil {
		return nil, err
	}

	return externalConfig, nil
}
