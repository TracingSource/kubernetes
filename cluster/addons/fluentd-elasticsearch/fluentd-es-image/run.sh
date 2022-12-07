#!/bin/sh

# These steps must be executed once the host /var and /lib volumes have
# been mounted, and therefore cannot be done in the docker build stage.

# For systems without journald
mkdir -p /var/log/journal

# Use exec to get the signal
# A non-quoted string and add the comment to prevent shellcheck failures on this line.
# See https://github.com/koalaman/shellcheck/wiki/SC2086
# shellcheck disable=SC2086
exec /usr/local/bin/fluentd $FLUENTD_ARGS
