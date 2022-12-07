package job

import (
	batch "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
)

func IsJobFinished(j *batch.Job) bool {
	for _, c := range j.Status.Conditions {
		if (c.Type == batch.JobComplete || c.Type == batch.JobFailed) && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
