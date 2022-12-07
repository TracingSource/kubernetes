package printers

import (
	"io"

	"k8s.io/apimachinery/pkg/runtime"
)

// NewDiscardingPrinter is a printer that discards all objects
func NewDiscardingPrinter() ResourcePrinterFunc {
	return ResourcePrinterFunc(func(runtime.Object, io.Writer) error {
		return nil
	})
}
