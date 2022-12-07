// deepcopy-gen is a tool for auto-generating DeepCopy functions.
//
// Given a list of input directories, it will generate functions that
// efficiently perform a full deep-copy of each type.  For any type that
// offers a `.DeepCopy()` method, it will simply call that.  Otherwise it will
// use standard value assignment whenever possible.  If that is not possible it
// will try to call its own generated copy function for the type, if the type is
// within the allowed root packages.  Failing that, it will fall back on
// `conversion.Cloner.DeepCopy(val)` to make the copy.  The resulting file will
// be stored in the same directory as the processed source package.
//
// Generation is governed by comment tags in the source.  Any package may
// request DeepCopy generation by including a comment in the file-comments of
// one file, of the form:
//   // +k8s:deepcopy-gen=package
//
// DeepCopy functions can be generated for individual types, rather than the
// entire package by specifying a comment on the type definion of the form:
//   // +k8s:deepcopy-gen=true
//
// When generating for a whole package, individual types may opt out of
// DeepCopy generation by specifying a comment on the of the form:
//   // +k8s:deepcopy-gen=false
//
// Note that registration is a whole-package option, and is not available for
// individual types.
package main

import (
	"flag"
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/deepcopy-gen/generators"
	"k8s.io/klog"

	generatorargs "k8s.io/code-generator/cmd/deepcopy-gen/args"
	"k8s.io/code-generator/pkg/util"
)

func main() {
	klog.InitFlags(nil)
	genericArgs, customArgs := generatorargs.NewDefaults()

	// Override defaults.
	// TODO: move this out of deepcopy-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())

	genericArgs.AddFlags(pflag.CommandLine)
	customArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := generatorargs.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	// Run it.
	if err := genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
