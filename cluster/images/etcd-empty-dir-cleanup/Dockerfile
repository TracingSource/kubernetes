FROM k8s.gcr.io/debian-base:v1.0.0

COPY etcdctl etcd-empty-dir-cleanup.sh /
RUN chmod a+rx /etcdctl /etcd-empty-dir-cleanup.sh

ENV ETCDCTL /etcdctl
ENV SLEEP_SECOND 3600

USER nobody:nogroup

ENTRYPOINT ["/etcd-empty-dir-cleanup.sh"]
