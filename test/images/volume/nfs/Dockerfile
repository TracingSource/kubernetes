FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/
RUN yum -y install /usr/bin/ps nfs-utils && yum clean all
RUN mkdir -p /exports
ADD run_nfs.sh /usr/local/bin/
ADD index.html /tmp/index.html
RUN chmod 644 /tmp/index.html

# mark /exports as a mount point
VOLUME /exports
# expose mountd 20048/tcp and nfsd 2049/tcp
EXPOSE 2049/tcp 20048/tcp

ENTRYPOINT ["/usr/local/bin/run_nfs.sh"]
CMD ["/exports", "/"]
