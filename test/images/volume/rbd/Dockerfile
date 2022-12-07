# CEPH all in one
# Based on image by Ricardo Rocha, ricardo@catalyst.net.nz

FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

# Base Packages
RUN yum install -y wget strace psmisc procps-ng ceph ceph-fuse && yum clean all

# Get ports exposed
EXPOSE 6789

ADD ./bootstrap.sh /bootstrap.sh
ADD ./mon.sh /mon.sh
ADD ./mds.sh /mds.sh
ADD ./osd.sh /osd.sh
ADD ./ceph.conf.sh /ceph.conf.sh
ADD ./keyring /var/lib/ceph/mon/keyring
ADD ./block.tar.gz /

CMD /bootstrap.sh
