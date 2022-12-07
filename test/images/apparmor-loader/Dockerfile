FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN apk add apparmor libapparmor --update-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing/ --allow-untrusted

ADD loader /usr/bin/loader

ENTRYPOINT ["/usr/bin/loader", "-logtostderr", "-v=2"]

# Default directory to watch.
CMD ["/profiles"]
