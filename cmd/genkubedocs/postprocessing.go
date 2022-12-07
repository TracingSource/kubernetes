package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// MarkdownPostProcessing goes though the generated files
func MarkdownPostProcessing(cmd *cobra.Command, dir string, processor func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := MarkdownPostProcessing(c, dir, processor); err != nil {
			return err
		}
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"
	filename := filepath.Join(dir, basename)

	markdownBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	processedMarkDown := processor(string(markdownBytes))

	return ioutil.WriteFile(filename, []byte(processedMarkDown), 0644)
}

// cleanupForInclude parts of markdown that will make difficult to use it as include in the website:
// - The title of the document (this allow more flexibility for include, e.g. include in tabs)
// - The sections see also, that assumes file will be used as a main page
func cleanupForInclude(md string) string {
	lines := strings.Split(md, "\n")

	cleanMd := ""
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if line == "### SEE ALSO" {
			break
		}

		cleanMd += line
		if i < len(lines)-1 {
			cleanMd += "\n"
		}
	}

	return cleanMd
}
