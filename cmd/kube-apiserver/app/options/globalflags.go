package options

import (
	"github.com/spf13/pflag"

	"k8s.io/component-base/cli/globalflag"

	// ensure libs have a chance to globally register their flags
	_ "k8s.io/apiserver/pkg/admission"
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
)

// AddCustomGlobalFlags explicitly registers flags that internal packages register
// against the global flagsets from "flag". We do this in order to prevent
// unwanted flags from leaking into the kube-apiserver's flagset.
func AddCustomGlobalFlags(fs *pflag.FlagSet) {
	// Lookup flags in global flag set and re-register the values with our flagset.

	// Adds flags from k8s.io/kubernetes/pkg/cloudprovider/providers.
	registerLegacyGlobalFlags(fs)

	// Adds flags from k8s.io/apiserver/pkg/admission.
	globalflag.Register(fs, "default-not-ready-toleration-seconds")
	globalflag.Register(fs, "default-unreachable-toleration-seconds")
}
