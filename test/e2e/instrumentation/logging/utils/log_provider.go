package utils

// LogProvider interface provides an API to get logs from the logging backend.
type LogProvider interface {
	Init() error
	Cleanup()
	ReadEntries(name string) []LogEntry
	LoggingAgentName() string
}
