package storage

import (
	"github.com/onsi/ginkgo"
	"k8s.io/kubernetes/test/e2e/storage/drivers"
	"k8s.io/kubernetes/test/e2e/storage/testsuites"
	"k8s.io/kubernetes/test/e2e/storage/utils"
)

// List of testDrivers to be executed in below loop
var testDrivers = []func() testsuites.TestDriver{
	drivers.InitNFSDriver,
	drivers.InitGlusterFSDriver,
	drivers.InitISCSIDriver,
	drivers.InitRbdDriver,
	drivers.InitCephFSDriver,
	drivers.InitHostPathDriver,
	drivers.InitHostPathSymlinkDriver,
	drivers.InitEmptydirDriver,
	drivers.InitCinderDriver,
	drivers.InitGcePdDriver,
	drivers.InitVSphereDriver,
	drivers.InitAzureDriver,
	drivers.InitAwsDriver,
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeDirectory),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeDirectoryLink),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeDirectoryBindMounted),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeDirectoryLinkBindMounted),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeTmpfs),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeBlock),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeBlockFS),
	drivers.InitLocalDriverWithVolumeType(utils.LocalVolumeGCELocalSSD),
}

// List of testSuites to be executed in below loop
var testSuites = []func() testsuites.TestSuite{
	testsuites.InitVolumesTestSuite,
	testsuites.InitVolumeIOTestSuite,
	testsuites.InitVolumeModeTestSuite,
	testsuites.InitSubPathTestSuite,
	testsuites.InitProvisioningTestSuite,
	testsuites.InitMultiVolumeTestSuite,
	testsuites.InitVolumeExpandTestSuite,
	testsuites.InitDisruptiveTestSuite,
	testsuites.InitVolumeLimitsTestSuite,
	testsuites.InitTopologyTestSuite,
}

// This executes testSuites for in-tree volumes.
var _ = utils.SIGDescribe("In-tree Volumes", func() {
	for _, initDriver := range testDrivers {
		curDriver := initDriver()

		ginkgo.Context(testsuites.GetDriverNameWithFeatureTags(curDriver), func() {
			testsuites.DefineTestSuite(curDriver, testSuites)
		})
	}
})
