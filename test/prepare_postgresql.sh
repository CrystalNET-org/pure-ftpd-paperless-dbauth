#!/bin/bash -x

psql_host=$(env | grep -ni "WP_SVC_.*DATABASE_PGSQL_SERVICE_HOST" | cut -d "=" -f 2)
echo "Waiting for postgresql to launch on ${psql_host}:5432..."

timeout 22 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' "${psql_host}" 5432

echo "postgresql launched, injecting sql"
psql -h "${psql_host}" -d test -p 5432 -U test -f ./test/prep_sql.sql