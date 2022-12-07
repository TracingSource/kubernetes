FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN apk add --no-cache redis

COPY redis.conf /etc/redis.conf

EXPOSE 6379
VOLUME /data

CMD ["redis-server", "/etc/redis.conf"]
