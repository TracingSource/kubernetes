// +build providerless

package options

import (
	"github.com/spf13/pflag"
)

func registerLegacyGlobalFlags(fs *pflag.FlagSet) {
	// no-op when no legacy providers are compiled in
}
