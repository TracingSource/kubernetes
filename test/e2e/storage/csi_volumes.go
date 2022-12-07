package storage

import (
	"k8s.io/kubernetes/test/e2e/storage/drivers"
	"k8s.io/kubernetes/test/e2e/storage/testsuites"
	"k8s.io/kubernetes/test/e2e/storage/utils"

	"github.com/onsi/ginkgo"
)

// List of testDrivers to be executed in below loop
var csiTestDrivers = []func() testsuites.TestDriver{
	drivers.InitHostPathCSIDriver,
	drivers.InitGcePDCSIDriver,
	// Don't run tests with mock driver (drivers.InitMockCSIDriver), it does not provide persistent storage.
}

// List of testSuites to be executed in below loop
var csiTestSuites = []func() testsuites.TestSuite{
	testsuites.InitEphemeralTestSuite,
	testsuites.InitVolumesTestSuite,
	testsuites.InitVolumeIOTestSuite,
	testsuites.InitVolumeModeTestSuite,
	testsuites.InitSubPathTestSuite,
	testsuites.InitProvisioningTestSuite,
	testsuites.InitSnapshottableTestSuite,
	testsuites.InitMultiVolumeTestSuite,
	testsuites.InitDisruptiveTestSuite,
	testsuites.InitVolumeExpandTestSuite,
	testsuites.InitVolumeLimitsTestSuite,
	testsuites.InitTopologyTestSuite,
}

// This executes testSuites for csi volumes.
var _ = utils.SIGDescribe("CSI Volumes", func() {
	for _, initDriver := range csiTestDrivers {
		curDriver := initDriver()

		ginkgo.Context(testsuites.GetDriverNameWithFeatureTags(curDriver), func() {
			testsuites.DefineTestSuite(curDriver, csiTestSuites)
		})
	}
})
