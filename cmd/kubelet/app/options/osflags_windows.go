// +build windows

package options

import (
	"github.com/spf13/pflag"
)

func (f *KubeletFlags) addOSFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&f.WindowsService, "windows-service", f.WindowsService, "Enable Windows Service Control Manager API integration")
}
