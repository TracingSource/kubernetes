FROM ubuntu:xenial

RUN apt-get update && apt-get install -y netbase nfs-common=1:1.2.8-9ubuntu12 glusterfs-client=3.7.6-1ubuntu1

ENTRYPOINT ["/bin/mount"]
