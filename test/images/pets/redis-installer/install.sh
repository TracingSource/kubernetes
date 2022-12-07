#!/bin/bash

# This volume is assumed to exist and is shared with parent of the init
# container. It contains the redis installation.
INSTALL_VOLUME="/opt"

# This volume is assumed to exist and is shared with the peer-finder
# init container. It contains on-start/change configuration scripts.
WORK_DIR="/work-dir"

VERSION="3.2.0"

for i in "$@"
do
case $i in
    -i=*|--install-into=*)
    INSTALL_VOLUME="${i#*=}"
    shift
    ;;
    -w=*|--work-dir=*)
    WORK_DIR="${i#*=}"
    shift
    ;;
    *)
    # unknown option
    ;;
esac
done

echo installing config scripts into "${WORK_DIR}"
mkdir -p "${WORK_DIR}"
cp /on-start.sh "${WORK_DIR}"/
cp /peer-finder "${WORK_DIR}"/

echo installing redis-"${VERSION}" into "${INSTALL_VOLUME}"
mkdir -p "${INSTALL_VOLUME}"
mv /redis "${INSTALL_VOLUME}"/redis
