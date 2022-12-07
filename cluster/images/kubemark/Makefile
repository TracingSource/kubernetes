# build Kubemark image from currently built binaries containing both 'real' master and Hollow Node.
# This makefile assumes that the kubemark binary is present in this directory.

# Allow the caller to override this.  Beware make's precedence. This:
#   REGISTRY=$VAR make
# .. is not the same as:
#   make REGISTRY=$VAR
REGISTRY := $(if $(REGISTRY),$(REGISTRY),staging-k8s.gcr.io)
IMAGE_TAG := $(if $(IMAGE_TAG),$(IMAGE_TAG),latest)

all: gcloudpush

build:
	docker build --pull -t $(REGISTRY)/kubemark:$(IMAGE_TAG) .

gcloudpush: build
	docker push $(REGISTRY)/kubemark:$(IMAGE_TAG)

push: build
	docker push $(REGISTRY)/kubemark:$(IMAGE_TAG)

.PHONY: all build gcloudpush push
