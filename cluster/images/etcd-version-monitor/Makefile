# Build the etcd-version-monitor image
#
# Usage:
# 	[GOLANG_VERSION=1.8.3] [REGISTRY=staging-k8s.gcr.io] [TAG=test] make (build|push)
# TODO(shyamjvs): Support architectures other than amd64 if needed.
ARCH:=amd64
GOLANG_VERSION?=1.8.3
REGISTRY?=staging-k8s.gcr.io
TAG?=0.1.3
IMAGE:=$(REGISTRY)/etcd-version-monitor:$(TAG)
CURRENT_DIR:=$(pwd)
TEMP_DIR:=$(shell mktemp -d)

build:
	# Copy the necessary files for building the image to TEMP_DIR.
	cp etcd-version-monitor.go Dockerfile $(TEMP_DIR)

	# Compile etcd-version-monitor.
	docker run -it \
	    -v $(shell pwd)/../../../:/go/src/k8s.io/kubernetes \
	    -v $(TEMP_DIR):/build \
	    -e GOARCH=$(ARCH) \
	    golang:$(GOLANG_VERSION) \
	    /bin/bash -c "CGO_ENABLED=0 go build -o /build/etcd-version-monitor k8s.io/kubernetes/cluster/images/etcd-version-monitor"

	docker build -t $(IMAGE) $(TEMP_DIR)

push: build
	docker push $(IMAGE)

all: build

.PHONY: build push
