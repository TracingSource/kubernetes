#!/usr/bin/env bash

DIR="$(mktemp -d)"

function start()
{
    mount -t tmpfs test "$DIR"
    chmod 755 "$DIR"
    cp /vol/* "$DIR/"
    /usr/sbin/glusterd -p /run/glusterd.pid
    gluster volume create test_vol "$(hostname -i):$DIR" force
    gluster volume start test_vol
}

function stop()
{
    gluster --mode=script volume stop test_vol force
    kill "$(cat /run/glusterd.pid)"
    umount "$DIR"
    rm -rf "$DIR"
    exit 0
}


trap stop TERM

start "$@"

while true; do
    sleep 5
done

