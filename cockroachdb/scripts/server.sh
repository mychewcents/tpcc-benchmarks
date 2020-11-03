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
  fi
else
  echo "Use the command \"cdbserv help\" to learn more about the acceptable parameters"
fi
