#!/bin/sh

echo "Removing empty directories from etcd..."

cleanup_empty_dirs () {
  if [ "$("${ETCDCTL}" ls "${1}")" ]; then
    for SUBDIR in $("${ETCDCTL}" ls -p "${1}" | grep "/$")
    do
      cleanup_empty_dirs "${SUBDIR}"
    done
  else
    echo "Removing empty key $1 ..."
    "${ETCDCTL}" rmdir "${1}"
  fi
}

while true
do
  echo "Starting cleanup..."
  cleanup_empty_dirs "/registry"
  echo "Done with cleanup."
  sleep "${SLEEP_SECOND}"
done
