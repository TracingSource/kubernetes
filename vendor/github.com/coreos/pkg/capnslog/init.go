// Copyright 2015 CoreOS, Inc.
//

//
// +build !windows

package capnslog

import (
	"io"
	"os"
	"syscall"
)

// Here's where the opinionation comes in. We need some sensible defaults,
// especially after taking over the log package. Your project (whatever it may
// be) may see things differently. That's okay; there should be no defaults in
// the main package that cannot be controlled or overridden programatically,
// otherwise it's a bug. Doing so is creating your own init_log.go file much
// like this one.

func init() {
	initHijack()

	// Go `log` pacakge uses os.Stderr.
	SetFormatter(NewDefaultFormatter(os.Stderr))
	SetGlobalLogLevel(INFO)
}

func NewDefaultFormatter(out io.Writer) Formatter {
	if syscall.Getppid() == 1 {
		// We're running under init, which may be systemd.
		f, err := NewJournaldFormatter()
		if err == nil {
			return f
		}
	}
	return NewPrettyFormatter(out, false)
}
