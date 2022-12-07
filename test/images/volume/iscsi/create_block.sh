#!/usr/bin/env bash

# Exit on the first error.
set -e

MNTDIR="$(mktemp -d)"

cleanup()
{
    # Make sure we return the right exit code
    RET=$?
    # Silently remove everything and ignore errors
    set +e
    /bin/umount "$MNTDIR" 2>/dev/null
    /bin/rmdir "$MNTDIR" 2>/dev/null
    /bin/rm block 2>/dev/null
    exit $RET
}

trap cleanup TERM EXIT

# Create 120MB device with ext2
# (volume_io tests need at least 100MB)
dd if=/dev/zero of=block seek=120 count=1 bs=1M
mkfs.ext2 block

# Add index.html to it
mount -o loop block "$MNTDIR"
echo "Hello from iscsi" > "$MNTDIR/index.html"
umount "$MNTDIR"

rm block.tar.gz 2>/dev/null || :
tar cfz block.tar.gz block
