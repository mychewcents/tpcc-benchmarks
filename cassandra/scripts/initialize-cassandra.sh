#!/bin/bash
cd "$(dirname "$0")/.."

rm -rf $CASSANDRA_HOME/data
rm -rf $CASSANDRA_HOME/logs

cqlsh --file 'scripts/cql/schema-initialization.cql'
cqlsh --file 'scripts/cql/data-initialization.cql'
