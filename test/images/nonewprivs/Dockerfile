FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

COPY nnp /usr/local/bin/nnp
RUN chmod +s /usr/local/bin/nnp

CMD ["nnp"]
