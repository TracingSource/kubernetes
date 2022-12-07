#!/usr/bin/env bash

#
# Configures /etc/ceph.conf from a template.
#

echo "
[global]
auth cluster required = none
auth service required = none
auth client required = none

[mon.a]
host = cephbox
mon addr = $1

[osd]
osd journal size = 128
journal dio = false

# allow running on ext4
osd max object name len = 256
osd max object namespace len = 64

[osd.0]
osd host = cephbox
" > /etc/ceph/ceph.conf

