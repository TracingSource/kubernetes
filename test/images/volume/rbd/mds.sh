#!/usr/bin/env bash

#
# Configures Ceph Metadata Service (mds), needed by CephFS
#
ceph-mds -i cephfs -c /etc/ceph/ceph.conf
