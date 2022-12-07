package logutil

import (
	"fmt"

	"github.com/coreos/pkg/capnslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var DefaultLogLevel = "info"

// ConvertToZapLevel converts log level string to zapcore.Level.
func ConvertToZapLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		panic(fmt.Sprintf("unknown level %q", lvl))
	}
}

// ConvertToCapnslogLogLevel convert log level string to capnslog.LogLevel.
// TODO: deprecate this in 3.5
func ConvertToCapnslogLogLevel(lvl string) capnslog.LogLevel {
	switch lvl {
	case "debug":
		return capnslog.DEBUG
	case "info":
		return capnslog.INFO
	case "warn":
		return capnslog.WARNING
	case "error":
		return capnslog.ERROR
	case "dpanic":
		return capnslog.CRITICAL
	case "panic":
		return capnslog.CRITICAL
	case "fatal":
		return capnslog.CRITICAL
	default:
		panic(fmt.Sprintf("unknown level %q", lvl))
	}
}
