#!/bin/bash

if [[ $1 == "help" ]]
then
  echo
  echo
  echo "This is simple manual for the executable."
  echo
  echo
  echo "\"cdbserv\" performs several functions:"
  echo
  echo
  echo "First Argument - action to perform"
  echo "   - init     : To initialize the Cockroach DB Instance"
  echo "   - start    : To start the CockroachDB server on the localhost"
  echo "   - stop     : To stop the local running instance of the Cockroach DB"
  echo "   - sql      : To start the SQL Client for the localhost"
  echo
  echo
  echo "Second Argument - host number to be used - Only used when creating directories"
  echo "   - host:port  : Location of the service of the cdb server instance"
  echo
  echo
  echo "Third Argument - host to connect to"
  echo "   - host:port  : Location of the service of the cdb server instance"
  echo
  echo
  echo "Happy Running!"
  echo
  echo
elif [[ $1 == "start" ]]
then
  if [[ "$#" -eq 5 ]]
  then
    cockroach start --insecure --store=$2 \
      --listen-addr=$3 \
      --http-addr=$4 \
      --join=$5 \
      --cache=.25 \
	    --max-sql-memory=.25 \
      --background
  else 
    echo "Incorrect number of arguments passed. Check \"help\" command."
  fi
elif [[ "$#" -eq 2 ]]
then
  if [[ $1 == "stop" ]]
  then
    cockroach quit --insecure --host=$2 --drain-wait 15s
  elif [[ $1 == "init" ]]
  then
    cockroach init --insecure --host=$2
  elif [[ $1 == "sql" ]]
  then
    cockroach sql --insecure --host=$2
  elif [[ $1 == "download-dataset" ]]
  then
    m -rf assets/data

    mkdir assets/data
    mkdir assets/data/raw
    mkdir assets/data/transactions
    mkdir assets/data/processed
    mkdir assets/data/processed/warehouse
    mkdir assets/data/processed/district
    mkdir assets/data/processed/customer
    mkdir assets/data/processed/order
    mkdir assets/data/processed/orderline
    mkdir assets/data/processed/stock
    mkdir assets/data/processed/item
    mkdir assets/data/processed/itempairs

    curl $2 -L -o assets/project-files.zip
    unzip assets/project-files.zip -d assets
    mv assets/project-files/data-files/* assets/data/raw
    mv assets/project-files/xact-files/* assets/data/transactions

    rm assets/project-files.zip
    rm -rf assets/project-files
  fi
elif [[ "$#" -eq 3 ]]
then
  if [[ $1 == 'setup-dirs' ]]
  then
    if [[ ! -d "${2}/cdb-server" ]]
      then
        echo "Creating the crdb server installation directory"
        mkdir -p $2/cdb-server
      fi

    extern_dir=$2/cdb-server/node-files/$3/extern 
    rm -rf $2/cdb-server/node-files/$3
    mkdir -p $extern_dir/assets/raw
    mkdir -p $extern_dir/assets/processed
    cp assets/data/raw/* $extern_dir/assets/raw
    cp -r assets/data/processed/* $extern_dir/assets/processed
  fi
else
  echo "Use the command \"cdbserv help\" to learn more about the acceptable parameters"
fi
