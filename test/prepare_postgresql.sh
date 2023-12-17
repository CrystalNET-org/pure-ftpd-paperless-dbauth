#!/bin/bash -x

psql_host=$(env | grep PORT_5432_TCP_ADDR | cut -d "=" -f 2)
echo "Waiting for postgresql to launch on ${psql_host}:5432..."

while ! nc -z "${psql_host}" 5432; do   
  sleep 0.1 # wait for 1/10 of the second before check again
done

echo "postgresql launched, injecting sql"
psql -h "${psql_host}" -d test -p 5432 -U test -f ./test/prep_psql.sql