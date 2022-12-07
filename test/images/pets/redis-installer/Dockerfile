# TODO: get rid of bash dependency and switch to plain busybox.
# The tar in busybox also doesn't seem to understand compression.
FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

# TODO: just use standard redis when there is one for 3.2.0.
RUN clean-install wget make gcc libc-dev bash

# See README.md
RUN wget -qO /redis-3.2.0.tar.gz http://download.redis.io/releases/redis-3.2.0.tar.gz && \
    tar -xzf /redis-3.2.0.tar.gz -C /tmp/ && rm /redis-3.2.0.tar.gz

# Clean out existing deps before installation
# see https://github.com/antirez/redis/issues/722
RUN cd /tmp/redis-3.2.0 && make distclean && mkdir -p /redis && \
    make install INSTALL_BIN=/redis && \
    mv /tmp/redis-3.2.0/redis.conf /redis/redis.conf && \
    rm -rf /tmp/redis-3.2.0

ADD on-start.sh /

COPY peer-finder /
ADD install.sh /
RUN chmod -c 755 /install.sh /on-start.sh /peer-finder
ENTRYPOINT ["/install.sh"]
