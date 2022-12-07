package spec

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// Debug is true when the SWAGGER_DEBUG env var is not empty.
	// It enables a more verbose logging of this package.
	Debug = os.Getenv("SWAGGER_DEBUG") != ""
	// specLogger is a debug logger for this package
	specLogger *log.Logger
)

func init() {
	debugOptions()
}

func debugOptions() {
	specLogger = log.New(os.Stdout, "spec:", log.LstdFlags)
}

func debugLog(msg string, args ...interface{}) {
	// A private, trivial trace logger, based on go-openapi/spec/expander.go:debugLog()
	if Debug {
		_, file1, pos1, _ := runtime.Caller(1)
		specLogger.Printf("%s:%d: %s", filepath.Base(file1), pos1, fmt.Sprintf(msg, args...))
	}
}
