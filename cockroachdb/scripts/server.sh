#! /bin/bash

if [ $1 == "help" ]
then
  echo
  echo
  echo "This is simple manual for the executable."
  echo
  echo
  echo "\"cdbserv\" accepts two command line arguments:"
  echo
  echo
  echo "First Argument - action to perform"
  echo "   - init  : To initialize the Cockroach DB Instance"
  echo "   - start : To start the CockroachDB server on the localhost"
  echo "   - stop  : To stop the local running instance of the Cockroach DB"
  echo "   - sql   : To start the SQL Client for the localhost"
  echo
  echo
  echo "Second Argument - host number to be used - Only used when creating directories"
  echo "   - <number>  : Should be [1, 5]"
  echo
  echo
  echo "Happy Running!"
  echo
  echo
elif [ "$#" -eq 2 ]
then
  if [ $1 == "start" ]
  then
    if [ $2 -gt 0 ] && [ $2 -lt 6 ]
    then
      printf "**********\nStarting Cockroach DB node on : node${2}\n**********\n\n"
      echo cockroach start --insecure --store=/home/stuproj/cs4224m/crdb-node-files/node$2 \
        --listen-addr=$(hostname -i):27000 \
        --http-addr=0.0.0.0:40000 \
        --join=192.168.48.179:27000,192.168.48.180:27000,192.168.48.181:27000,192.168.48.182:27000,192.168.48.183:27000 \
        --background

      printf "\n**********\nStarted server on : node${2}\n**********\n"
    else
      echo "Use the host number from 1 to 5 ONLY...."
    fi
  elif [ $1 == "stop" ]
  then
    printf "**********\nStopping the node : node${2}\n**********\n\n"
    cockroach quit --insecure --host=$(hostname -i):27000
    printf "\n**********\nStopped the node : node${2}\n**********\n"
  elif [ $1 == "init" ]
  then
    printf "**********\nInitializing the cluster : node${2}\n**********\n\n"
    cockroach init --insecure --host=$(hostname -i):27000
    printf "\n**********\nStopped the node : node${2}\n**********\n"
  elif [ $1 == "sql" ]
  then
    printf "**********\nStarting the SQL Client : node${2}\n***********\n\n"
    cockroach sql --insecure --host=$(hostname -i):27000
    printf "\n**********\nStopped the SQL Client : node${2}\n**********\n\n"
  fi
else
  echo "Use the command \"cdbserv help\" to learn more about the acceptable parameters"
fi
