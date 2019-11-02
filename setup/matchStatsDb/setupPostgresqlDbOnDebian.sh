#!/usr/bin/env bash

# This script sets up a new postgres cluster and exposes it on all interfaces
# It assumes a debian server, you may need to adapt it for other distros

PG_VERSION=11
PG_CLUSTER=ssl
PG_DB=ssl_match_stats

# create a new cluster
runuser -l postgres -c pg_createcluster $PG_VERSION $PG_CLUSTER

# start new cluster
systemctl daemon-reload
systemctl start postgresql@$PG_VERSION-$PG_CLUSTER

# Allow access to database '$PG_DB' for all users and all IPs
echo "host $PG_DB all 0.0.0.0/0 md5" >> /etc/postgresql/$PG_VERSION/$PG_CLUSTER/pg_hba.conf

# Listen on all IPs of the server
echo "listen_addresses = '*'" >> /etc/postgresql/$PG_VERSION/$PG_CLUSTER/postgresql.conf
