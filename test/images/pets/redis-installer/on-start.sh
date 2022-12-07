#!/usr/bin/env bash

set -e

CFG=/opt/redis/redis.conf
HOSTNAME=$(hostname)
DATADIR="/data"
# Port on which redis listens for connections.
PORT=6379

# Ping everyone but ourself to see if there's a master. Only one pet starts at
# a time, so if we don't see a master we can assume the position is ours.
while read -ra LINE; do
    if [[ "${LINE[0]}" == *"${HOSTNAME}"* ]]; then
        sed -i -e "s|^bind.*$|bind ${LINE[0]}|" ${CFG}
    elif [ "$(/opt/redis/redis-cli -h "${LINE[0]}" info | grep role | sed 's,\r$,,')" = "role:master" ]; then
        # TODO: More restrictive regex?
        sed -i -e "s|^# slaveof.*$|slaveof ${LINE[0]} ${PORT}|" ${CFG}
    fi
done

# Set the data directory for append only log and snapshot files. This should
# be a persistent volume for consistency.
sed -i -e "s|^.*dir .*$|dir ${DATADIR}|" ${CFG}

# The append only log is written for every SET operation. Without this setting,
# redis just snapshots periodically which is only safe for a cache. This will
# produce an appendonly.aof file in the configured data dir.
sed -i -e "s|^appendonly .*$|appendonly yes|" ${CFG}

# Every write triggers an fsync. Recommended default is "everysec", which
# is only safe for AP applications.
sed -i -e "s|^appendfsync .*$|appendfsync always|" ${CFG}


