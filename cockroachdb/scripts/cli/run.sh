#! /bin/bash

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
  echo "   - local  : To run a single client instance"
  echo "   - prod   : To run parallel client instances"
  echo
  echo
  echo "Second Argument - Experiment number to run"
  echo "   - <number>    : Should be [5, 8]"
  echo
  echo
  echo "Third Argument - Host to run on or the client instance to run with"
  echo "   - <number>    : - Acts as a host number when run with \"prod\""
  echo "                   - Acts as a client number when run with \"local\""
  echo
  echo
  echo "Happy Running!"
  echo
  echo
elif [ "$#" -eq 4 ] 
then
  if [ $1 == "dev" ] || [ $1 == "prod" ]
  then
    env=$1
    exp=$2
    config_file_path=$4
    if [ $exp -gt 4 ] && [ $exp -lt 9 ] 
    then
      if [ $env == 'dev' ]
      then
        client=$3
        ./clientCmd -exp=$exp -client=$client -config=$config_file_path -node=1 < assets/data/transactions/$client.txt
      elif [ $env == 'prod' ]
      then
        node=$3
        skip=5
        total=20

        case $exp in
          5) skip=4
          ;;
          6) skip=5
          ;;
          7) skip=4
          total=40
          ;;
          8) skip=5
          total=40
          ;;
        esac

        for (( i=node; i<=total; i=i+skip ))
        do
          ./clientCmd -exp=$exp -client=$i -config=$config_file_path -node=$node < assets/data/transactions/$i.txt &
        done
      fi
    fi
  fi
else
  echo "Use the \"help\" command to learn how to use this file"
fi
