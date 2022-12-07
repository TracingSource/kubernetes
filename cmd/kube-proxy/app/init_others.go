// +build !windows

package app

import (
	"github.com/spf13/pflag"
)

func initForOS(service bool) error {
	return nil
}

func (o *Options) addOSFlags(fs *pflag.FlagSet) {
}
