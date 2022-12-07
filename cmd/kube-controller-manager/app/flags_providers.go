// +build !providerless

package app

import (
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
)

func registerLegacyGlobalFlags(namedFlagSets cliflag.NamedFlagSets) {
	// hoist this flag from the global flagset to preserve the commandline until
	// the gce cloudprovider is removed.
	globalflag.Register(namedFlagSets.FlagSet("generic"), "cloud-provider-gce-lb-src-cidrs")
	namedFlagSets.FlagSet("generic").MarkDeprecated("cloud-provider-gce-lb-src-cidrs", "This flag will be removed once the GCE Cloud Provider is removed from kube-controller-manager")
}
