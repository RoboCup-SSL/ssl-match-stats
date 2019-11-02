#!/usr/bin/env bash

# Create new database for the statistics data
# This assumes a debian system and root permissions

DB_CLUSTER=11/ssl
DB_NAME=ssl_match_stats
ADMIN_USER=ssl_match_stats
VIEW_USER=ssl_match_stats_view
ADMIN_PW=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
VIEW_PW=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

echo "PW for admin user: $ADMIN_PW"
echo "PW for view user: $VIEW_PW"

cat << EOF > /tmp/matchStats.sql
CREATE DATABASE $DB_NAME;

CREATE USER $ADMIN_USER WITH ENCRYPTED PASSWORD '$ADMIN_PW';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $ADMIN_USER;

CREATE USER $VIEW_USER WITH ENCRYPTED PASSWORD '$VIEW_PW';
GRANT CONNECT ON DATABASE $DB_NAME TO $VIEW_USER;
GRANT USAGE ON SCHEMA public TO $VIEW_USER;
-- SELECT must be assigned to individual tables
EOF

runuser -l postgres -c "psql --cluster $DB_CLUSTER -f /tmp/matchStats.sql"
rm /tmp/matchStats.sql
