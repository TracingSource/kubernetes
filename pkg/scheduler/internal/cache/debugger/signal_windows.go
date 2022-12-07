package debugger

import "os"

// compareSignal is the signal to trigger cache compare. For windows,
// it's SIGINT.
var compareSignal = os.Interrupt
