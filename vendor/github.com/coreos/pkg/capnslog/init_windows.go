package capnslog

import "os"

func init() {
	initHijack()

	// Go `log` package uses os.Stderr.
	SetFormatter(NewPrettyFormatter(os.Stderr, false))
	SetGlobalLogLevel(INFO)
}
