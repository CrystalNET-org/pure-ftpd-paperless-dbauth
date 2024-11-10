#!/bin/bash -x

mysql_host=$(env | grep "WP_SVC_.*DATABASE_MYSQL_SERVICE_HOST" | cut -d "=" -f 2)
echo "Waiting for mysql to launch on ${mysql_host}:3306..."

timeout 22 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' "${mysql_host}" 3306

echo "mysql launched, injecting sql"
mysql --host="${mysql_host}" --port=3306 --user=root --password=example --database=test < ./test/prep_sql.sql