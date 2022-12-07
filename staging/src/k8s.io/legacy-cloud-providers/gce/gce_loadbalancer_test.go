// +build !providerless

package gce

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLoadBalancer(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService("")

	// When a loadbalancer has not been created
	status, found, err := gce.GetLoadBalancer(context.Background(), vals.ClusterName, apiService)
	assert.Nil(t, status)
	assert.False(t, found)
	assert.Nil(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)
	expectedStatus, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	require.NoError(t, err)

	status, found, err = gce.GetLoadBalancer(context.Background(), vals.ClusterName, apiService)
	assert.Equal(t, expectedStatus, status)
	assert.True(t, found)
	assert.Nil(t, err)
}

func TestEnsureLoadBalancerCreatesExternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService("")
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	assert.NoError(t, err)
	assert.NotEmpty(t, status.Ingress)
	assertExternalLbResources(t, gce, apiService, vals, nodeNames)
}

func TestEnsureLoadBalancerCreatesInternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService(string(LBTypeInternal))
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	assert.NoError(t, err)
	assert.NotEmpty(t, status.Ingress)
	assertInternalLbResources(t, gce, apiService, vals, nodeNames)
}

func TestEnsureLoadBalancerDeletesExistingInternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService("")
	createInternalLoadBalancer(gce, apiService, nil, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	assert.NoError(t, err)
	assert.NotEmpty(t, status.Ingress)

	assertExternalLbResources(t, gce, apiService, vals, nodeNames)
	assertInternalLbResourcesDeleted(t, gce, apiService, vals, false)
}

func TestEnsureLoadBalancerDeletesExistingExternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	nodes, err := createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService("")
	createExternalLoadBalancer(gce, apiService, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	apiService = fakeLoadbalancerService(string(LBTypeInternal))
	status, err := gce.EnsureLoadBalancer(context.Background(), vals.ClusterName, apiService, nodes)
	assert.NoError(t, err)
	assert.NotEmpty(t, status.Ingress)

	assertInternalLbResources(t, gce, apiService, vals, nodeNames)
	assertExternalLbResourcesDeleted(t, gce, apiService, vals, false)
}

func TestEnsureLoadBalancerDeletedDeletesExternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	_, err = createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService("")
	createExternalLoadBalancer(gce, apiService, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	err = gce.EnsureLoadBalancerDeleted(context.Background(), vals.ClusterName, apiService)
	assert.NoError(t, err)
	assertExternalLbResourcesDeleted(t, gce, apiService, vals, true)
}

func TestEnsureLoadBalancerDeletedDeletesInternalLb(t *testing.T) {
	t.Parallel()

	vals := DefaultTestClusterValues()
	gce, err := fakeGCECloud(vals)
	require.NoError(t, err)

	nodeNames := []string{"test-node-1"}
	_, err = createAndInsertNodes(gce, nodeNames, vals.ZoneName)
	require.NoError(t, err)

	apiService := fakeLoadbalancerService(string(LBTypeInternal))
	createInternalLoadBalancer(gce, apiService, nil, nodeNames, vals.ClusterName, vals.ClusterID, vals.ZoneName)

	err = gce.EnsureLoadBalancerDeleted(context.Background(), vals.ClusterName, apiService)
	assert.NoError(t, err)
	assertInternalLbResourcesDeleted(t, gce, apiService, vals, true)
}
