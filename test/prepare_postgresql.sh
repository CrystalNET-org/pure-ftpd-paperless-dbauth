#!/bin/bash

sleep 30s
psql_host=$(env | grep PORT_5432_TCP_ADDR | cut -d "=" -f 2)
echo "connecting to ${psql_host}"
psql -h "${psql_host}" -d test -p 5432 -U test -f ./sql/prep_psql.sql