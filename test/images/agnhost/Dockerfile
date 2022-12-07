FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

# from hostexec image
# install necessary packages:
# - curl, nc: used by a lot of e2e tests
# - iproute2: includes ss used in NodePort tests
# from iperf image
# install necessary packages: iperf, bash
RUN apk --update add bind-tools curl netcat-openbsd iproute2 iperf bash && rm -rf /var/cache/apk/* \
  && ln -s /usr/bin/iperf /usr/local/bin/iperf \
  && ls -altrh /usr/local/bin/iperf

# PORT 8080 needed by: netexec, nettest
# PORT 8081 needed by: netexec
# PORT 9376 needed by: serve-hostname
EXPOSE 8080 8081 9376

# from netexec
RUN mkdir /uploads

# from porter
ADD porter/localhost.crt localhost.crt
ADD porter/localhost.key localhost.key

ADD agnhost agnhost

# needed for the entrypoint-tester related tests. Some of the entrypoint-tester related tests
# overrides this image's entrypoint with agnhost-2 binary, and will verify that the correct
# entrypoint is used by the containers.
RUN ln -s agnhost agnhost-2

ENTRYPOINT ["/agnhost"]
CMD ["pause"]
