package options

import (
	"flag"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
)

func TestAddCustomGlobalFlags(t *testing.T) {
	namedFlagSets := &cliflag.NamedFlagSets{}

	// Note that we will register all flags (including klog flags) into the same
	// flag set. This allows us to test against all global flags from
	// flags.CommandLine.
	nfs := namedFlagSets.FlagSet("test")
	globalflag.AddGlobalFlags(nfs, "test-cmd")
	AddCustomGlobalFlags(nfs)

	actualFlag := []string{}
	nfs.VisitAll(func(flag *pflag.Flag) {
		actualFlag = append(actualFlag, flag.Name)
	})

	// Get all flags from flags.CommandLine, except flag `test.*`.
	wantedFlag := []string{"help"}
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.VisitAll(func(flag *pflag.Flag) {
		if !strings.Contains(flag.Name, "test.") {
			wantedFlag = append(wantedFlag, flag.Name)
		}
	})
	sort.Strings(wantedFlag)

	if !reflect.DeepEqual(wantedFlag, actualFlag) {
		t.Errorf("[Default]: expected %+v, got %+v", wantedFlag, actualFlag)
	}
}
