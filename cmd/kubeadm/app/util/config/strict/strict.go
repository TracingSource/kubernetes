package strict

import (
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	"k8s.io/kubernetes/cmd/kubeadm/app/componentconfigs"
	"sigs.k8s.io/yaml"
)

// VerifyUnmarshalStrict takes a YAML byte slice and a GroupVersionKind and verifies if the YAML
// schema is known and if it unmarshals with strict mode.
//
// TODO(neolit123): The returned error here is currently ignored everywhere and a klog warning is thrown instead.
// We don't want to turn this into an actual error yet. Eventually this can be controlled with an optional CLI flag.
func VerifyUnmarshalStrict(bytes []byte, gvk schema.GroupVersionKind) error {

	var (
		iface interface{}
		err   error
	)

	iface, err = scheme.Scheme.New(gvk)
	if err != nil {
		iface, err = componentconfigs.Scheme.New(gvk)
		if err != nil {
			err := errors.Errorf("unknown configuration %#v for scheme definitions in %q and %q",
				gvk, scheme.Scheme.Name(), componentconfigs.Scheme.Name())
			klog.Warning(err.Error())
			return err
		}
	}

	if err := yaml.UnmarshalStrict(bytes, iface); err != nil {
		err := errors.Wrapf(err, "error unmarshaling configuration %#v", gvk)
		klog.Warning(err.Error())
		return err
	}
	return nil
}
