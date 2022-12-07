package system

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// ValidationResultType is type of the validation result. Different validation results
// corresponds to different colors.
type ValidationResultType int32

const (
	good ValidationResultType = iota
	bad
	warn
)

// color is the color of the message.
type color int32

const (
	red    color = 31
	green  color = 32
	yellow color = 33
	white  color = 37
)

func colorize(s string, c color) string {
	return fmt.Sprintf("\033[0;%dm%s\033[0m", c, s)
}

// StreamReporter is the default reporter for the system verification test.
type StreamReporter struct {
	// The stream that this reporter is writing to
	WriteStream io.Writer
}

// Report reports validation result in different color depending on the result type.
func (dr *StreamReporter) Report(key, value string, resultType ValidationResultType) error {
	var c color
	switch resultType {
	case good:
		c = green
	case bad:
		c = red
	case warn:
		c = yellow
	default:
		c = white
	}
	if dr.WriteStream == nil {
		return errors.New("WriteStream has to be defined for this reporter")
	}

	fmt.Fprintf(dr.WriteStream, "%s: %s\n", colorize(key, white), colorize(value, c))
	return nil
}

// DefaultReporter is the default Reporter
var DefaultReporter = &StreamReporter{
	WriteStream: os.Stdout,
}
