# Build the hyperkube image.
#
# Usage:
#   [ARCH=amd64] [REGISTRY="staging-k8s.gcr.io"] make (build|push) VERSION={some_released_version_of_kubernetes}

REGISTRY?=staging-k8s.gcr.io
ARCH?=amd64
OUT_DIR?=_output

BASEIMAGE=k8s.gcr.io/debian-hyperkube-base-$(ARCH):0.12.1

LOCAL_OUTPUT_PATH=$(shell pwd)/../../../$(OUT_DIR)/local/bin/linux/$(ARCH)
DOCKERIZED_OUTPUT_PATH=$(shell pwd)/../../../$(OUT_DIR)/dockerized/bin/linux/$(ARCH)
OUTPUT_PATH?=$(shell test -d $(LOCAL_OUTPUT_PATH) && echo $(LOCAL_OUTPUT_PATH) || echo $(DOCKERIZED_OUTPUT_PATH))

TEMP_DIR:=$(shell mktemp -d -t hyperkubeXXXXXX)

all: build

build:

ifndef VERSION
    $(error VERSION is undefined)
endif
	cp -r ./* ${TEMP_DIR}

	tar -cvzf ${TEMP_DIR}/binaries.tgz  -C ${OUTPUT_PATH} kube-apiserver kube-controller-manager \
		kube-proxy kube-scheduler kubectl kubelet

	chmod a+rx ${TEMP_DIR}/hyperkube

	cd ${TEMP_DIR} && sed -i.back "s|BASEIMAGE|${BASEIMAGE}|g" Dockerfile

	docker build --pull -t ${REGISTRY}/hyperkube-${ARCH}:${VERSION} ${TEMP_DIR}
	rm -rf "${TEMP_DIR}"

push: build
	docker push ${REGISTRY}/hyperkube-${ARCH}:${VERSION}
ifeq ($(ARCH),amd64)
	docker rmi ${REGISTRY}/hyperkube:${VERSION} 2>/dev/null || true
	docker tag ${REGISTRY}/hyperkube-${ARCH}:${VERSION} ${REGISTRY}/hyperkube:${VERSION}
	docker push ${REGISTRY}/hyperkube:${VERSION}
endif

.PHONY: build push all
