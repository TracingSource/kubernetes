package main

import (
	"flag"
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/klog"

	generatorargs "k8s.io/code-generator/cmd/register-gen/args"
	"k8s.io/code-generator/cmd/register-gen/generators"
	"k8s.io/code-generator/pkg/util"
	"k8s.io/gengo/args"
)

func main() {
	klog.InitFlags(nil)
	genericArgs := generatorargs.NewDefaults()
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), util.BoilerplatePath())
	genericArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()
	if err := generatorargs.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	if err := genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
