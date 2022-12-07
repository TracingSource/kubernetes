/*
Copyright 2016 Euan Kemp


*/

package kmsgparser

import stdlog "log"

// Logger is a glog compatible logging interface
// The StandardLogger struct can be used to wrap a log.Logger from the golang
// "log" package to create a standard a logger fulfilling this interface as
// well.
type Logger interface {
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// StandardLogger adapts the "log" package's Logger interface to be a Logger
type StandardLogger struct {
	*stdlog.Logger
}

func (s *StandardLogger) Warningf(fmt string, args ...interface{}) {
	if s.Logger == nil {
		return
	}
	s.Logger.Printf("[WARNING] "+fmt, args)
}

func (s *StandardLogger) Infof(fmt string, args ...interface{}) {
	if s.Logger == nil {
		return
	}
	s.Logger.Printf("[INFO] "+fmt, args)
}

func (s *StandardLogger) Errorf(fmt string, args ...interface{}) {
	if s.Logger == nil {
		return
	}
	s.Logger.Printf("[INFO] "+fmt, args)
}
