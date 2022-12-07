FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN clean-install wget bash dnsutils

COPY peer-finder /

EXPOSE 9376
ENTRYPOINT ["/peer-finder"]
