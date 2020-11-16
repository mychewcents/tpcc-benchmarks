#!/bin/bash

cockroach sql --insecure --host=$1 -e "${2}" --format=csv > $3