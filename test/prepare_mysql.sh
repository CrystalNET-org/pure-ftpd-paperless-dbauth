#!/bin/bash -x

mysql_host=$(env | grep PORT_3306_TCP_ADDR | cut -d "=" -f 2)
echo "Waiting for mysql to launch on ${mysql_host}:3306..."

timeout 22 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' "${mysql_host}" 3306

echo "mysql launched, injecting sql"
mysql --host="${mysql_host}" --port=5432 --user=test --password=example --database=test < ./test/prep_sql.sql