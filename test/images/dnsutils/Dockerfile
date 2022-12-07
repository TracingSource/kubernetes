FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN apk add --no-cache bind-tools dnsmasq

ADD https://github.com/coredns/coredns/releases/download/v1.5.0/coredns_1.5.0_linux_BASEARCH.tgz /coredns.tgz
RUN tar -xzvf /coredns.tgz && rm -f /coredns.tgz
