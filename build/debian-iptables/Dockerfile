FROM BASEIMAGE

# Install latest iptables package from buster-backports
RUN echo deb http://deb.debian.org/debian buster-backports main >> /etc/apt/sources.list; \
    apt-get update; \
    apt-get -t buster-backports -y --no-install-recommends install iptables

# Install other dependencies and then clean up apt caches
RUN clean-install \
    conntrack \
    ebtables \
    ipset \
    kmod \
    netbase

# Install iptables wrapper scripts to detect the correct iptables mode
# the first time any of them is run
COPY iptables-wrapper /usr/sbin/iptables-wrapper

RUN update-alternatives \
	--install /usr/sbin/iptables iptables /usr/sbin/iptables-wrapper 100 \
	--slave /usr/sbin/iptables-restore iptables-restore /usr/sbin/iptables-wrapper \
	--slave /usr/sbin/iptables-save iptables-save /usr/sbin/iptables-wrapper
RUN update-alternatives \
	--install /usr/sbin/ip6tables ip6tables /usr/sbin/iptables-wrapper 100 \
	--slave /usr/sbin/ip6tables-restore ip6tables-restore /usr/sbin/iptables-wrapper \
	--slave /usr/sbin/ip6tables-save ip6tables-save /usr/sbin/iptables-wrapper
