#!/bin/bash

cqlsh --file 'scripts/cql/schema-initialization.cql'
cqlsh --file 'scripts/cql/data-initialization.cql'
