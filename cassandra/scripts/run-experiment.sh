#!/bin/bash
cd "$(dirname ${BASH_SOURCE[0]})/.."

if [ $1 == "help" ]
then
  echo
  echo
  echo "This is a simple manual for the \"run.sh\" file."
  echo
  echo
  echo "\"run.sh\" accepts 3 command line argument and SHOULD BE run from the \"cockroachdb\" directory"
  echo
  echo
  echo "First Argument - Type of the environment"
  echo "   - dev    : To run a single client instance"
  echo "   - prod   : To run parallel client instances"
  echo
  echo
  echo "Second Argument - Experiment number to run"
  echo "   - <number>    : Should be [5, 8]"
  echo
  echo
  echo "Third Argument - Host to run on or the client instance to run with"
  echo "   - <number>    : - Acts as a host number when run with \"prod\""
  echo "                   - Acts as a client number when run with \"dev\""
  echo
  echo
  echo "Happy Running!"
  echo
  echo
elif [ "$#" -eq 3 ] 
then
  if [ $1 == "local" ] || [ $1 == "prod" ]
  then
    env=$1
    exp=$2
    host_or_client=$3
    if [ $exp -gt 0 ] && [ $exp -lt 5 ]
    then
      if [ $env == 'local' ]
      then
        go run cmd/client/cassandra-client.go $exp $host_or_client configs/$env/cassandra-config.xml < assets/data/transactions/0.txt &
      elif [ $env == 'prod' ]
      then
        skip=5
        total=20

        case $exp in
          3) total=40
          ;;
          4) total=40
          ;;
        esac

        for (( i=host_or_client; i<=total; i=i+skip ))
        do
          go run cmd/client/cassandra-client.go $exp $i configs/$env/experiment-$exp-cassandra-config.xml < assets/data/transactions/$i.txt &
        done
      fi
    fi
  fi
else
  echo "Use the \"help\" command to learn how to use this file"
fi
