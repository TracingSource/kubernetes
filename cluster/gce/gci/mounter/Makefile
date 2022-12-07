TAG=v2
REGISTRY=staging-k8s.gcr.io
IMAGE=gci-mounter

all: container

container:
	docker build --pull -t ${REGISTRY}/${IMAGE}:${TAG} .

push:
	docker push ${REGISTRY}/${IMAGE}:${TAG}

upload:
	./stage-upload.sh ${TAG} ${REGISTRY}/${IMAGE}:${TAG}

.PHONY: all container push
