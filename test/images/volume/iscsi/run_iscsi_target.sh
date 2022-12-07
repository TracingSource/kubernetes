#!/usr/bin/env bash

# This script does not run any daemon, it only configures iSCSI target (=server)
# in kernel. It is possible to run this script multiple times on a single
# node, each run will create its own IQN and LUN.

# Kubernetes must provide unique name.
IQN=$1

# targetcli synchronizes over dbus, however it does not work in
# containers. Use flock instead
LOCK=/srv/iscsi/targetcli.lock

function start()
{
    # targetcli need dbus. It may not run on the host, so start a private one
    mkdir /run/dbus
    dbus-daemon  --system

    # Create new IQN (iSCSI Qualified Name)
    flock $LOCK targetcli /iscsi create "$IQN"
    # Run it in demo mode, i.e. no authentication
    flock $LOCK targetcli /iscsi/"$IQN"/tpg1 set attribute authentication=0 demo_mode_write_protect=0 generate_node_acls=1 cache_dynamic_acls=1

    # Create unique "block volume" (i.e. flat file) on the *host*.
    # Having it in the container confuses kernel from some reason
    # and it's not able to server multiple LUNs from different
    # containers.
    # /srv/iscsi must be bind-mount from the host.
    cp /block /srv/iscsi/"$IQN"

    # Make the block volume available through our IQN as LUN 0
    flock $LOCK targetcli /backstores/fileio create block-"$IQN" /srv/iscsi/"$IQN"
    flock $LOCK targetcli /iscsi/"$IQN"/tpg1/luns create /backstores/fileio/block-"$IQN"

    echo "iscsi target started"
}

function stop()
{
    echo "stopping iscsi target"
    # Remove IQN
    flock $LOCK targetcli /iscsi/"$IQN"/tpg1/luns/ delete 0
    flock $LOCK targetcli /iscsi delete "$IQN"
    # Remove block device mapping
    flock $LOCK targetcli /backstores/fileio delete block-"$IQN"
    /bin/rm -f /srv/iscsi/"$IQN"
    echo "iscsi target stopped"
    exit 0
}


trap stop TERM
start

while true; do
    sleep 1
done
