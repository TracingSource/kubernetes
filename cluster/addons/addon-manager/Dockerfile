FROM BASEIMAGE

RUN clean-install bash

ADD kube-addons.sh /opt/
ADD kubectl /usr/local/bin/

CMD ["/opt/kube-addons.sh"]
