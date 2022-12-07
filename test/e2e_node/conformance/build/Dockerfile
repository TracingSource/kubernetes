FROM BASEIMAGE

COPY ginkgo /usr/local/bin/
COPY e2e_node.test /usr/local/bin

# This is a placeholder that will be substituted in the Makefile.
COPY_SYSTEM_SPEC_FILE

# The following environment variables can be override when starting the container.
# FOCUS is regex matching test to run. By default run all conformance test.
# SKIP is regex matching test to skip. By default skip flaky and serial test.
# PARALLELISM is the number of processes the test will run in parallel.
# REPORT_PATH is the path in the container to save test result and logs.
# FLAKE_ATTEMPTS is the time to retry when there is a test failure. By default 2.
# TEST_ARGS is the test arguments passed into the test.
# EXTRA_ENVS is the extra environment variables needed for node e2e tests.
ENV FOCUS="\[Conformance\]" \
	   SKIP="\[Flaky\]|\[Serial\]" \
	   PARALLELISM=8 \
	   REPORT_PATH="/var/result" \
	   FLAKE_ATTEMPTS=2 \
	   TEST_ARGS="" \
	   EXTRA_ENVS=""

ENTRYPOINT ginkgo --focus="$FOCUS" \
	--skip="$SKIP" \
	--nodes=$PARALLELISM \
	--flakeAttempts=$FLAKE_ATTEMPTS \
	/usr/local/bin/e2e_node.test \
	-- --conformance=true \
	--prepull-images=false \
	--report-dir="$REPORT_PATH" \
	# This is a placeholder that will be substituted in the Makefile.
	--system-spec-name=SYSTEM_SPEC_NAME \
	# This is a placeholder that will be substituted in the Makefile.
	--system-spec-file=SYSTEM_SPEC_FILE_PATH \
	--extra-envs=$EXTRA_ENVS \
	$TEST_ARGS
