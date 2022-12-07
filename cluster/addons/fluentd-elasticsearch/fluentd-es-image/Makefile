.PHONY:	build push

PREFIX = quay.io/fluentd_elasticsearch
IMAGE = fluentd
TAG = v2.7.0

build:
	docker build --tag ${PREFIX}/${IMAGE}:${TAG} .
	docker build --tag ${PREFIX}/${IMAGE}:latest .

push:
	docker push ${PREFIX}/${IMAGE}:${TAG}
	docker push ${PREFIX}/${IMAGE}:latest
