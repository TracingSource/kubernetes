package args

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/defaulter-gen/generators"
)

// CustomArgs is used by the gengo framework to pass args specific to this generator.
type CustomArgs generators.CustomArgs

// NewDefaults returns default arguments for the generator.
func NewDefaults() (*args.GeneratorArgs, *CustomArgs) {
	genericArgs := args.Default().WithoutDefaultFlagParsing()
	customArgs := &CustomArgs{}
	genericArgs.CustomArgs = (*generators.CustomArgs)(customArgs) // convert to upstream type to make type-casts work there
	genericArgs.OutputFileBaseName = "zz_generated.defaults"
	return genericArgs, customArgs
}

// AddFlags add the generator flags to the flag set.
func (ca *CustomArgs) AddFlags(fs *pflag.FlagSet) {
	pflag.CommandLine.StringSliceVar(&ca.ExtraPeerDirs, "extra-peer-dirs", ca.ExtraPeerDirs,
		"Comma-separated list of import paths which are considered, after tag-specified peers, for conversions.")
}

// Validate checks the given arguments.
func Validate(genericArgs *args.GeneratorArgs) error {
	_ = genericArgs.CustomArgs.(*generators.CustomArgs)

	if len(genericArgs.OutputFileBaseName) == 0 {
		return fmt.Errorf("output file base name cannot be empty")
	}

	return nil
}
