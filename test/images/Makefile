include  ../../hack/make-rules/Makefile.manifest

REGISTRY ?= gcr.io/kubernetes-e2e-test-images
GOARM ?= 7
QEMUVERSION=v2.9.1
GOLANG_VERSION=1.13.5
export

ifndef WHAT
$(error WHAT is a required variable, ex: make all WHAT=net)
endif

# Build code.
#
# Args:
#   WHAT: Directory names to build.
#
# Example:
#   make all WHAT=clusterapi-tester
all: all-container

all-container:
	./image-util.sh build $(WHAT)

all-push: all-container
	./image-util.sh push $(WHAT)

.PHONY: all all-push all-container
