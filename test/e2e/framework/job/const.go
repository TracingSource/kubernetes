package job

import "time"

const (
	// JobTimeout is how long to wait for a job to finish.
	JobTimeout = 15 * time.Minute

	// JobSelectorKey is a job selector name
	JobSelectorKey = "job"
)
