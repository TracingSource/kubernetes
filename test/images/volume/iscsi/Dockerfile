FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN yum install -y targetcli && yum clean all
ADD run_iscsi_target.sh /usr/local/bin/
ADD block.tar.gz /

ENTRYPOINT ["/usr/local/bin/run_iscsi_target.sh"]
