package genericclioptions

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/printers"
)

// KubeTemplatePrintFlags composes print flags that provide both a JSONPath and a go-template printer.
// This is necessary if dealing with cases that require support both both printers, since both sets of flags
// require overlapping flags.
type KubeTemplatePrintFlags struct {
	GoTemplatePrintFlags *GoTemplatePrintFlags
	JSONPathPrintFlags   *JSONPathPrintFlags

	AllowMissingKeys *bool
	TemplateArgument *string
}

func (f *KubeTemplatePrintFlags) AllowedFormats() []string {
	if f == nil {
		return []string{}
	}
	return append(f.GoTemplatePrintFlags.AllowedFormats(), f.JSONPathPrintFlags.AllowedFormats()...)
}

func (f *KubeTemplatePrintFlags) ToPrinter(outputFormat string) (printers.ResourcePrinter, error) {
	if f == nil {
		return nil, NoCompatiblePrinterError{}
	}

	if p, err := f.JSONPathPrintFlags.ToPrinter(outputFormat); !IsNoCompatiblePrinterError(err) {
		return p, err
	}
	return f.GoTemplatePrintFlags.ToPrinter(outputFormat)
}

// AddFlags receives a *cobra.Command reference and binds
// flags related to template printing to it
func (f *KubeTemplatePrintFlags) AddFlags(c *cobra.Command) {
	if f == nil {
		return
	}

	if f.TemplateArgument != nil {
		c.Flags().StringVar(f.TemplateArgument, "template", *f.TemplateArgument, "Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].")
		c.MarkFlagFilename("template")
	}
	if f.AllowMissingKeys != nil {
		c.Flags().BoolVar(f.AllowMissingKeys, "allow-missing-template-keys", *f.AllowMissingKeys, "If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats.")
	}
}

// NewKubeTemplatePrintFlags returns flags associated with
// --template printing, with default values set.
func NewKubeTemplatePrintFlags() *KubeTemplatePrintFlags {
	allowMissingKeysPtr := true
	templateArgPtr := ""

	return &KubeTemplatePrintFlags{
		GoTemplatePrintFlags: &GoTemplatePrintFlags{
			TemplateArgument: &templateArgPtr,
			AllowMissingKeys: &allowMissingKeysPtr,
		},
		JSONPathPrintFlags: &JSONPathPrintFlags{
			TemplateArgument: &templateArgPtr,
			AllowMissingKeys: &allowMissingKeysPtr,
		},

		TemplateArgument: &templateArgPtr,
		AllowMissingKeys: &allowMissingKeysPtr,
	}
}
