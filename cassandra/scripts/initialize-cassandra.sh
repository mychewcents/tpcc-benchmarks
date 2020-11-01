#!/bin/bash
cd "$(dirname ${BASH_SOURCE[0]})/.."
cqlsh --file 'scripts/cql/schema-initialization.cql'
cqlsh --file 'scripts/cql/data-initialization.cql'
