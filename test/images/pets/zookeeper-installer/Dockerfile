# TODO: get rid of bash dependency and switch to plain busybox.
# The tar in busybox also doesn't seem to understand compression.
FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN clean-install wget bash netcat

ADD on-start.sh /

COPY peer-finder /

# See README.md
RUN wget -q -O /zookeeper-3.5.0-alpha.tar.gz http://archive.apache.org/dist/zookeeper/zookeeper-3.5.0-alpha/zookeeper-3.5.0-alpha.tar.gz && \
    tar -xzf /zookeeper-3.5.0-alpha.tar.gz -C /tmp/ && mv /tmp/zookeeper-3.5.0-alpha /zookeeper && rm /zookeeper-3.5.0-alpha.tar.gz

ADD install.sh /
RUN chmod -c 755 /install.sh /on-start.sh /peer-finder
ENTRYPOINT ["/install.sh"]
