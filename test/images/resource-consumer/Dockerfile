FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN clean-install stress

ADD consumer /consumer
ADD consume-cpu /consume-cpu
EXPOSE 8080
ENTRYPOINT ["/consumer"]
