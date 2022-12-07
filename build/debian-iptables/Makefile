.PHONY:	build push all all-build all-push-images all-push push-manifest

REGISTRY?="staging-k8s.gcr.io"
IMAGE=$(REGISTRY)/debian-iptables
TAG?=v12.0.1
ARCH?=amd64
ALL_ARCH = amd64 arm arm64 ppc64le s390x
TEMP_DIR:=$(shell mktemp -d)

BASEIMAGE?=k8s.gcr.io/debian-base-$(ARCH):v2.0.0

# This option is for running docker manifest command
export DOCKER_CLI_EXPERIMENTAL := enabled

SUDO=$(if $(filter 0,$(shell id -u)),,sudo)

build:
	cp ./* $(TEMP_DIR)
	cd $(TEMP_DIR) && sed -i "s|BASEIMAGE|$(BASEIMAGE)|g" Dockerfile

ifneq ($(ARCH),amd64)
	# Register /usr/bin/qemu-ARCH-static as the handler for non-x86 binaries in the kernel
	$(SUDO) ../../third_party/multiarch/qemu-user-static/register/register.sh --reset
endif

	docker build --pull -t $(IMAGE)-$(ARCH):$(TAG) $(TEMP_DIR)

push: build
	docker push $(IMAGE)-$(ARCH):$(TAG)

sub-build-%:
	$(MAKE) ARCH=$* build

all-build: $(addprefix sub-build-,$(ALL_ARCH))

sub-push-image-%:
	$(MAKE) ARCH=$* push

all-push-images: $(addprefix sub-push-image-,$(ALL_ARCH))

all-push: all-push-images push-manifest

push-manifest:
	docker manifest create --amend $(IMAGE):$(TAG) $(shell echo $(ALL_ARCH) | sed -e "s~[^ ]*~$(IMAGE)\-&:$(TAG)~g")
	@for arch in $(ALL_ARCH); do docker manifest annotate --arch $${arch} ${IMAGE}:${TAG} ${IMAGE}-$${arch}:${TAG}; done
	docker manifest push --purge ${IMAGE}:${TAG}

all: all-push
