#!/usr/bin/env bash

#
# Configures and launches a new MON.
#

# monitor setup
monmaptool --create --clobber --fsid "$(uuidgen)" --add a "${1}":6789 /etc/ceph/monmap
mkdir /var/lib/ceph/mon/ceph-a
ceph-mon -i a --mkfs --monmap /etc/ceph/monmap -k /var/lib/ceph/mon/keyring
cp /var/lib/ceph/mon/keyring /var/lib/ceph/mon/ceph-a
ceph-mon -i a --monmap /etc/ceph/monmap -k /var/lib/ceph/mon/ceph-a/keyring

# client setup (handy)
cp /var/lib/ceph/mon/keyring /etc/ceph

# for this test we want to
ceph osd getcrushmap -o /tmp/crushc
crushtool -d /tmp/crushc -o /tmp/crushd
sed -i 's/step chooseleaf firstn 0 type host/step chooseleaf firstn 0 type osd/' /tmp/crushd
crushtool -c /tmp/crushd -o /tmp/crushc
ceph osd setcrushmap -i /tmp/crushc
