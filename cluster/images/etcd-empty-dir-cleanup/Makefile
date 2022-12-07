.PHONY:	build push

ETCD_VERSION = 3.4.3
# Image should be pulled from k8s.gcr.io, which will auto-detect
# region (us, eu, asia, ...) and pull from the closest.
REGISTRY = k8s.gcr.io
# Images should be pushed to staging-k8s.gcr.io.
PUSH_REGISTRY = staging-k8s.gcr.io
TAG = 3.4.3.0

clean:
	rm -rf etcdctl etcd-v$(ETCD_VERSION)-linux-amd64 etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz

build: clean
	curl -L -O https://github.com/coreos/etcd/releases/download/v$(ETCD_VERSION)/etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz
	tar xzvf etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz
	cp etcd-v$(ETCD_VERSION)-linux-amd64/etcdctl .
	docker build --pull -t $(REGISTRY)/etcd-empty-dir-cleanup:$(TAG) .
	rm -rf etcdctl etcd-v$(ETCD_VERSION)-linux-amd64 etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz

push: build
	docker tag $(REGISTRY)/etcd-empty-dir-cleanup:$(TAG) $(PUSH_REGISTRY)/etcd-empty-dir-cleanup:$(TAG)
	docker push $(PUSH_REGISTRY)/etcd-empty-dir-cleanup:$(TAG)
