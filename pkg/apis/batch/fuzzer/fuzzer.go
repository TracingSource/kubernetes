package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/batch"
)

func newBool(val bool) *bool {
	p := new(bool)
	*p = val
	return p
}

// Funcs returns the fuzzer functions for the batch api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(j *batch.Job, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again

			// match defaulting
			if len(j.Labels) == 0 {
				j.Labels = j.Spec.Template.Labels
			}
		},
		func(j *batch.JobSpec, c fuzz.Continue) {
			c.FuzzNoCustom(j) // fuzz self without calling this function again
			completions := int32(c.Rand.Int31())
			parallelism := int32(c.Rand.Int31())
			backoffLimit := int32(c.Rand.Int31())
			j.Completions = &completions
			j.Parallelism = &parallelism
			j.BackoffLimit = &backoffLimit
			if c.Rand.Int31()%2 == 0 {
				j.ManualSelector = newBool(true)
			} else {
				j.ManualSelector = nil
			}
		},
		func(sj *batch.CronJobSpec, c fuzz.Continue) {
			c.FuzzNoCustom(sj)
			suspend := c.RandBool()
			sj.Suspend = &suspend
			sds := int64(c.RandUint64())
			sj.StartingDeadlineSeconds = &sds
			sj.Schedule = c.RandString()
			successfulJobsHistoryLimit := int32(c.Rand.Int31())
			sj.SuccessfulJobsHistoryLimit = &successfulJobsHistoryLimit
			failedJobsHistoryLimit := int32(c.Rand.Int31())
			sj.FailedJobsHistoryLimit = &failedJobsHistoryLimit
		},
		func(cp *batch.ConcurrencyPolicy, c fuzz.Continue) {
			policies := []batch.ConcurrencyPolicy{batch.AllowConcurrent, batch.ForbidConcurrent, batch.ReplaceConcurrent}
			*cp = policies[c.Rand.Intn(len(policies))]
		},
	}
}
