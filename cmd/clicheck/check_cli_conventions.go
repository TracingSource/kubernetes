package main

import (
	"fmt"
	"io/ioutil"
	"os"

	cmdsanity "k8s.io/kubectl/pkg/cmd/util/sanity"
	"k8s.io/kubernetes/pkg/kubectl/cmd"
)

func main() {
	var errorCount int

	kubectl := cmd.NewKubectlCommand(os.Stdin, ioutil.Discard, ioutil.Discard)
	errors := cmdsanity.RunCmdChecks(kubectl, cmdsanity.AllCmdChecks, []string{})
	for _, err := range errors {
		errorCount++
		fmt.Fprintf(os.Stderr, "     %d. %s\n", errorCount, err)
	}

	errors = cmdsanity.RunGlobalChecks(cmdsanity.AllGlobalChecks)
	for _, err := range errors {
		errorCount++
		fmt.Fprintf(os.Stderr, "     %d. %s\n", errorCount, err)
	}

	if errorCount > 0 {
		fmt.Fprintf(os.Stdout, "Found %d errors.\n", errorCount)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "Congrats, CLI looks good!")
}
