#!/bin/bash

mysql_host=$(env | grep "WP_SVC_.*DATABASE_MYSQL_SERVICE_HOST" | cut -d "=" -f 2)
echo "Waiting for mysql to launch on ${mysql_host}:3306..."

timeout 22 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' "${mysql_host}" 3306

export PAPERLESS_DBHOST=${mysql_host}
export PAPERLESS_DBPORT="3306"
export PAPERLESS_DBNAME="test"
export DB_USER="root"
export PAPERLESS_DBPASS="example"
export PAPERLESS_DBENGINE="mysql"  # or "mysql" for MariaDB
export PAPERLESS_CONSUMPTION_DIR="/directory"

export AUTHD_ACCOUNT=test
export AUTHD_PASSWORD=testpassword

DEBUG=true out/verify_pw_amd64