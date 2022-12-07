// Package upgrades provides a framework for testing Kubernetes
// features before, during, and after different types of upgrades.
package upgrades

import (
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/kubernetes/test/e2e/framework"
)

// UpgradeType represents different types of upgrades.
type UpgradeType int

const (
	// MasterUpgrade indicates that only the master is being upgraded.
	MasterUpgrade UpgradeType = iota

	// NodeUpgrade indicates that only the nodes are being upgraded.
	NodeUpgrade

	// ClusterUpgrade indicates that both master and nodes are
	// being upgraded.
	ClusterUpgrade

	// EtcdUpgrade indicates that only etcd is being upgraded (or migrated
	// between storage versions).
	EtcdUpgrade
)

// Test is an interface for upgrade tests.
type Test interface {
	// Name should return a test name sans spaces.
	Name() string

	// Setup should create and verify whatever objects need to
	// exist before the upgrade disruption starts.
	Setup(f *framework.Framework)

	// Test will run during the upgrade. When the upgrade is
	// complete, done will be closed and final validation can
	// begin.
	Test(f *framework.Framework, done <-chan struct{}, upgrade UpgradeType)

	// Teardown should clean up any objects that are created that
	// aren't already cleaned up by the framework. This will
	// always be called, even if Setup failed.
	Teardown(f *framework.Framework)
}

// Skippable is an interface that an upgrade test can implement to be
// able to indicate that it should be skipped.
type Skippable interface {
	// Skip should return true if test should be skipped. upgCtx
	// provides information about the upgrade that is going to
	// occur.
	Skip(upgCtx UpgradeContext) bool
}

// UpgradeContext contains information about all the stages of the
// upgrade that is going to occur.
type UpgradeContext struct {
	Versions []VersionContext
}

// VersionContext represents a stage of the upgrade.
type VersionContext struct {
	Version   version.Version
	NodeImage string
}
