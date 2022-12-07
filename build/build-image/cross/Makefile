.PHONY:	build push

REGISTRY?=staging-k8s.gcr.io
IMAGE=$(REGISTRY)/kube-cross
TAG=$(shell cat VERSION)


all: push

build:
	docker build --pull -t $(IMAGE):$(TAG) .

push: build
	docker push $(IMAGE):$(TAG)
