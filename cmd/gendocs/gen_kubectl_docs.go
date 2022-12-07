package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra/doc"
	"k8s.io/kubernetes/cmd/genutils"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/"
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory]\n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutils.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by kubectl so the output is consistent,
	// regardless of where we run.
	os.Setenv("HOME", "/home/username")
	// TODO os.Stdin should really be something like ioutil.Discard, but a Reader
	kubectl := cmd.NewKubectlCommand(os.Stdin, ioutil.Discard, ioutil.Discard)
	doc.GenMarkdownTree(kubectl, outDir)
}
