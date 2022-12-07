package main

const (
	// resultsTarballName is the name of the tarball we create with all the results.
	resultsTarballName = "e2e.tar.gz"

	// doneFileName is the name of the file that signals to the Sonobuoy worker we are
	// done. The file should contain the path to the results file.
	doneFileName = "done"

	// resultsDirEnvKey is the env var which stores which directory to put the donefile
	// and results into. It is a shared, mounted volume between the plugin and Sonobuoy.
	resultsDirEnvKey = "RESULTS_DIR"

	// logFileName is the name of the file which stdout is tee'd to.
	logFileName = "e2e.log"

	// Misc env vars which were explicitly supported prior to the go runner.
	dryRunEnvKey     = "E2E_DRYRUN"
	parallelEnvKey   = "E2E_PARALLEL"
	focusEnvKey      = "E2E_FOCUS"
	skipEnvKey       = "E2E_SKIP"
	providerEnvKey   = "E2E_PROVIDER"
	kubeconfigEnvKey = "KUBECONFIG"
	ginkgoEnvKey     = "GINKGO_BIN"
	testBinEnvKey    = "TEST_BIN"

	// extraGinkgoArgsEnvKey, if set, will is a list of other arguments to pass to ginkgo.
	// These are passed before the test binary and include things like `--afterSuiteHook`.
	extraGinkgoArgsEnvKey = "E2E_EXTRA_GINKGO_ARGS"

	// extraArgsEnvKey, if set, will is a list of other arguments to pass to the tests.
	// These are passed after the `--` and include things like `--provider`.
	extraArgsEnvKey = "E2E_EXTRA_ARGS"

	// extraArgsSeparaterEnvKey specifies how to split the extra args values. If unset,
	// it will default to splitting by spaces.
	extraArgsSeparaterEnvKey = "E2E_EXTRA_ARGS_SEP"

	defaultSkip         = ""
	defaultFocus        = "[Conformance]"
	defaultProvider     = "local"
	defaultParallel     = "1"
	defaultResultsDir   = "/tmp/results"
	defaultGinkgoBinary = "/usr/local/bin/ginkgo"
	defaultTestBinary   = "/usr/local/bin/e2e.test"

	// serialTestsRegexp is the default skip value if running in parallel. Will not
	// override an explicit E2E_SKIP value.
	serialTestsRegexp = "[Serial]"
)
