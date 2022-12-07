FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

# WARNING: Please note that the script below removes the security packages from arm64 and ppc64el images
# as they do not exist anymore in the debian repositories for jessie. So we do not recommend using this
# image for any production use and limit use of this image to just test scenarios.

COPY fixup-apt-list.sh /
RUN ["/fixup-apt-list.sh"]

RUN apt-get -q update && \
    apt-get install -y dnsutils && \
    apt-get clean

ADD https://github.com/coredns/coredns/releases/download/v1.5.0/coredns_1.5.0_linux_BASEARCH.tgz /coredns.tgz
RUN tar -xzvf /coredns.tgz && rm -f /coredns.tgz
