package main

import (
	"flag"
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/code-generator/cmd/lister-gen/generators"
	"k8s.io/code-generator/pkg/util"
	"k8s.io/gengo/args"
	"k8s.io/klog"

	generatorargs "k8s.io/code-generator/cmd/lister-gen/args"
)

func main() {
	klog.InitFlags(nil)
	genericArgs, customArgs := generatorargs.NewDefaults()

	// Override defaults.
	// TODO: move this out of lister-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())
	genericArgs.OutputPackagePath = "k8s.io/kubernetes/pkg/client/listers"

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
