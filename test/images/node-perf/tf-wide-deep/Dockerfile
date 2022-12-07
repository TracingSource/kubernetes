FROM BASEIMAGE

CROSS_BUILD_COPY qemu-QEMUARCH-static /usr/bin/

RUN apt-get update && apt-get install -y wget time

RUN case $(uname -m) in \
    aarch64) \
      pip install tensorflow-aarch64; \
      ;; \
    *) \
      pip install tensorflow; \
      ;; \
    esac

RUN wget https://github.com/tensorflow/models/archive/v1.9.0.tar.gz \
&& tar xzf v1.9.0.tar.gz \
&& rm -f v1.9.0.tar.gz

WORKDIR $HOME/models-1.9.0/official/wide_deep
ENV PYTHONPATH $PYTHONPATH:$HOME/models-1.9.0

ENTRYPOINT python ./wide_deep.py
