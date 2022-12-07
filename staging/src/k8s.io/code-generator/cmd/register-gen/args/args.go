package args

import (
	"fmt"

	"k8s.io/gengo/args"
)

// NewDefaults returns default arguments for the generator.
func NewDefaults() *args.GeneratorArgs {
	genericArgs := args.Default().WithoutDefaultFlagParsing()
	genericArgs.OutputFileBaseName = "zz_generated.register"
	return genericArgs
}

// Validate checks the given arguments.
func Validate(genericArgs *args.GeneratorArgs) error {
	if len(genericArgs.OutputFileBaseName) == 0 {
		return fmt.Errorf("output file base name cannot be empty")
	}

	return nil
}
