package eventratelimit

import (
	"time"
)

// realClock implements flowcontrol.Clock in terms of standard time functions.
type realClock struct{}

// Now is identical to time.Now.
func (realClock) Now() time.Time {
	return time.Now()
}

// Sleep is identical to time.Sleep.
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
