// +build providerless

package app

import (
	cliflag "k8s.io/component-base/cli/flag"
)

func registerLegacyGlobalFlags(namedFlagSets cliflag.NamedFlagSets) {
	// no-op when legacy cloud providers are not compiled
}
