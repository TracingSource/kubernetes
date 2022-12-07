// +build !providerless

package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAzureMetricLabelCardinality(t *testing.T) {
	mc := newMetricContext("test", "create", "resource_group", "subscription_id", "source")
	assert.Len(t, mc.attributes, len(metricLabels), "cardinalities of labels and values must match")
}

func TestAzureMetricLabelPrefix(t *testing.T) {
	mc := newMetricContext("prefix", "request", "resource_group", "subscription_id", "source")
	found := false
	for _, attribute := range mc.attributes {
		if attribute == "prefix_request" {
			found = true
		}
	}
	assert.True(t, found, "request label must be prefixed")
}
