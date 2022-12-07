SRCS=regression-issue-74839
ARCH ?= amd64
TARGET ?= $(CURDIR)
GOARM ?= 7
GOLANG_VERSION ?= latest
SRC_DIR = $(notdir $(shell pwd))

export

.PHONY: bin
bin:
	../image-util.sh bin $(SRCS)

