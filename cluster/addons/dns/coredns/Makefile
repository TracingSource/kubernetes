# Makefile for the kubedns underscore templates to Salt/Pillar and other formats.

# If you update the *.base templates, please run this Makefile before pushing.
#
# Usage:
#    make

all: transform

# .base -> .in pattern rule
%.in: %.base
	sed -f transforms2salt.sed $< | sed s/__SOURCE_FILENAME__/$</g > $@

# .base -> .sed pattern rule
%.sed: %.base
	sed -f transforms2sed.sed $<  | sed s/__SOURCE_FILENAME__/$</g > $@

transform: coredns.yaml.in coredns.yaml.sed

.PHONY: transform
