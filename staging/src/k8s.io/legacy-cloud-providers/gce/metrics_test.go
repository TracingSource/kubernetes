// +build !providerless

package gce

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyMetricLabelCardinality(t *testing.T) {
	mc := newGenericMetricContext("foo", "get", "us-central1", "<n/a>", "alpha")
	assert.Len(t, mc.attributes, len(metricLabels), "cardinalities of labels and values must match")
}
