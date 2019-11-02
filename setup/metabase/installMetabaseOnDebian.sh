#!/bin/bash

# Based on: https://www.metabase.com/docs/latest/operations-guide/running-metabase-on-debian.html

set -e

INSTALL_DIR=/srv/metabase
DB_PW=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
mkdir "$INSTALL_DIR"
wget -O "$INSTALL_DIR/metabase.jar" http://downloads.metabase.com/v0.33.4/metabase.jar

groupadd -r metabase
useradd -r -s /bin/false -g metabase metabase
chown -R metabase:metabase "$INSTALL_DIR"
touch /var/log/metabase.log
chown metabase:metabase /var/log/metabase.log
touch /etc/default/metabase
chmod 640 /etc/default/metabase

cat << EOF > /etc/systemd/system/metabase.service
[Unit]
Description=Metabase server
After=syslog.target
After=network.target
After=postgresql.target

[Service]
WorkingDirectory=$INSTALL_DIR
ExecStart=/usr/bin/java -jar $INSTALL_DIR/metabase.jar
EnvironmentFile=/etc/default/metabase
User=metabase
Type=simple
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=metabase
SuccessExitStatus=143
TimeoutStopSec=120
Restart=always

[Install]
WantedBy=multi-user.target
EOF

cat << EOF > /etc/rsyslog.d/metabase.conf
if $programname == 'metabase' then /var/log/metabase.log
& stop
EOF

systemctl restart rsyslog.service


cat << EOF > /etc/default/metabase
MB_PASSWORD_COMPLEXITY=normal
MB_PASSWORD_LENGTH=10
MB_JETTY_HOST=0.0.0.0
MB_JETTY_PORT=8083
MB_DB_TYPE=postgres
MB_DB_DBNAME=metabase
MB_DB_PORT=5432
MB_DB_USER=metabase
MB_DB_PASS=$DB_PW
MB_DB_HOST=localhost
MB_EMOJI_IN_LOGS=true
EOF

cat << EOF > /tmp/metabase.sql
CREATE DATABASE metabase;
CREATE USER metabase WITH ENCRYPTED PASSWORD '$DB_PW';
GRANT ALL PRIVILEGES ON DATABASE metabase TO metabase;
EOF

runuser -l postgres -c 'psql -f /tmp/metabase.sql'
rm /tmp/metabase.sql

systemctl daemon-reload
systemctl enable metabase
systemctl start metabase
systemctl status metabase


echo "Mostly done. You only need to add metabase to your nginx proxy"
