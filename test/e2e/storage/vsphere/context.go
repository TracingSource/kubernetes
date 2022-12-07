package vsphere

// Context holds common information for vSphere tests
type VSphereContext struct {
	NodeMapper       *NodeMapper
	VSphereInstances map[string]*VSphere
}

// TestContext should be used by all tests to access common context data. It should be initialized only once, during bootstrapping the tests.
var TestContext VSphereContext
