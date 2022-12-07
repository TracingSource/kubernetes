FROM BASEIMAGE

# TODO(#69896): deprecate the shortened aliases in /
RUN ln -s /hyperkube /apiserver \
 && ln -s /hyperkube /cloud-controller-manager \
 && ln -s /hyperkube /controller-manager \
 && ln -s /hyperkube /kubectl \
 && ln -s /hyperkube /kubelet \
 && ln -s /hyperkube /proxy \
 && ln -s /hyperkube /scheduler \
 && ln -s /hyperkube /usr/local/bin/cloud-controller-manager \
 && ln -s /hyperkube /usr/local/bin/kube-apiserver \
 && ln -s /hyperkube /usr/local/bin/kube-controller-manager \
 && ln -s /hyperkube /usr/local/bin/kube-proxy \
 && ln -s /hyperkube /usr/local/bin/kube-scheduler \
 && ln -s /hyperkube /usr/local/bin/kubectl \
 && ln -s /hyperkube /usr/local/bin/kubelet

RUN echo CACHEBUST>/dev/null && clean-install \
    bash

# The samba-common, cifs-utils, and nfs-common packages depend on
# ucf, which itself depends on /bin/bash.
RUN echo "dash dash/sh boolean false" | debconf-set-selections
RUN DEBIAN_FRONTEND=noninteractive dpkg-reconfigure dash

RUN echo CACHEBUST>/dev/null && clean-install \
    ca-certificates \
    ceph-common \
    cifs-utils \
    conntrack \
    e2fsprogs \
    xfsprogs \
    ebtables \
    ethtool \
    git \
    glusterfs-client \
    iptables \
    ipset \
    jq \
    kmod \
    openssh-client \
    netbase \
    nfs-common \
    socat \
    udev \
    util-linux

COPY cni-bin/bin /opt/cni/bin
