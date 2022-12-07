package logreduction

import (
	"sync"
	"time"
)

var nowfunc = func() time.Time { return time.Now() }

// LogReduction provides a filter for consecutive identical log messages;
// a message will be printed no more than once per interval.
// If a string of messages is interrupted by a different message,
// the interval timer will be reset.
type LogReduction struct {
	lastError           map[string]string
	errorPrinted        map[string]time.Time
	errorMapLock        sync.Mutex
	identicalErrorDelay time.Duration
}

// NewLogReduction returns an initialized LogReduction
func NewLogReduction(identicalErrorDelay time.Duration) *LogReduction {
	l := new(LogReduction)
	l.lastError = make(map[string]string)
	l.errorPrinted = make(map[string]time.Time)
	l.identicalErrorDelay = identicalErrorDelay
	return l
}

func (l *LogReduction) cleanupErrorTimeouts() {
	for name, timeout := range l.errorPrinted {
		if nowfunc().Sub(timeout) >= l.identicalErrorDelay {
			delete(l.errorPrinted, name)
			delete(l.lastError, name)
		}
	}
}

// ShouldMessageBePrinted determines whether a message should be printed based
// on how long ago this particular message was last printed
func (l *LogReduction) ShouldMessageBePrinted(message string, parentID string) bool {
	l.errorMapLock.Lock()
	defer l.errorMapLock.Unlock()
	l.cleanupErrorTimeouts()
	lastMsg, ok := l.lastError[parentID]
	lastPrinted, ok1 := l.errorPrinted[parentID]
	if !ok || !ok1 || message != lastMsg || nowfunc().Sub(lastPrinted) >= l.identicalErrorDelay {
		l.errorPrinted[parentID] = nowfunc()
		l.lastError[parentID] = message
		return true
	}
	return false
}

// ClearID clears out log reduction records pertaining to a particular parent
// (e. g. container ID)
func (l *LogReduction) ClearID(parentID string) {
	l.errorMapLock.Lock()
	defer l.errorMapLock.Unlock()
	delete(l.lastError, parentID)
	delete(l.errorPrinted, parentID)
}
