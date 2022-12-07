FROM BASEIMAGE

# Cleanup all the soft links
ADD binaries.tgz /usr/local/bin

# Copy the shell script
COPY hyperkube /hyperkube

ENTRYPOINT ["/hyperkube"]