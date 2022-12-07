SRCS = consumer consume-cpu/consume-cpu
ARCH ?= amd64
TARGET ?= $(CURDIR)
GOLANG_VERSION ?= latest
SRC_DIR = $(notdir $(shell pwd))
export

bin:
	../image-util.sh bin $(SRCS)

.PHONY: bin
