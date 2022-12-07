FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN yum -y install hostname glusterfs-server && yum clean all
ADD glusterd.vol /etc/glusterfs/
ADD run_gluster.sh /usr/local/bin/
ADD index.html /vol/
RUN chmod 644 /vol/index.html

EXPOSE 24007/tcp 49152/tcp

ENTRYPOINT ["/usr/local/bin/run_gluster.sh"]
