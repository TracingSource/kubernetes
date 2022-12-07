.PHONY:	build push

PREFIX = quay.io/fluentd_elasticsearch
IMAGE = elasticsearch
TAG = v7.3.2

build:
	docker build --tag ${PREFIX}/${IMAGE}:${TAG} .
	docker build --tag ${PREFIX}/${IMAGE}:latest .

push:
	docker push ${PREFIX}/${IMAGE}:${TAG}
	docker push ${PREFIX}/${IMAGE}:latest
