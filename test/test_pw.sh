#!/bin/bash

psql_host=$(env | grep PORT_5432_TCP_ADDR | cut -d "=" -f 2)
echo "Waiting for postgresql to launch on ${psql_host}:5432..."

timeout 22 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' "${psql_host}" 5432

export PAPERLESS_DBHOST=${psql_host}
export PAPERLESS_DBPORT="5432"
export PAPERLESS_DBNAME=${POSTGRES_DB}
export DB_USER="test"
export PAPERLESS_DBPASS="example"
export PAPERLESS_DBENGINE="postgres"  # or "mysql" for MariaDB
export PAPERLESS_CONSUMPTION_DIR="/directory"

export AUTHD_ACCOUNT=test
export AUTHD_PASSWORD=testpassword

DEBUG=true out/verify_pw_amd64