// +build !windows

package options

import (
	"github.com/spf13/pflag"
)

func (f *KubeletFlags) addOSFlags(fs *pflag.FlagSet) {
}
