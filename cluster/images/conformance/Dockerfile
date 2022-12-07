FROM BASEIMAGE

COPY ginkgo /usr/local/bin/
COPY e2e.test /usr/local/bin/
COPY kubectl /usr/local/bin/
COPY run_e2e.sh /run_e2e.sh
COPY gorunner /gorunner
COPY cluster /kubernetes/cluster
WORKDIR /usr/local/bin

ENV E2E_FOCUS="\[Conformance\]"
ENV E2E_SKIP=""
ENV E2E_PROVIDER="local"
ENV E2E_PARALLEL="1"
ENV E2E_VERBOSITY="4"
ENV RESULTS_DIR="/tmp/results"
ENV KUBECONFIG=""

CMD [ "/bin/bash", "-c", "/run_e2e.sh" ]
