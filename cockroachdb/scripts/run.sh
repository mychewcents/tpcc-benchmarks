#! /bin/sh

env=$1
exp=$2
host_or_client=$3

if [ "$#" -eq 3 ] 
then
  if [ $exp -gt 4 ] && [ $exp -lt 9 ] 
  then
    go build -o app cmd/app/main.go
    chmod a+x ./app

    if [ $env == 'dev' ]
    then
      ./app -exp=$exp -client=$host_or_client < assets/data/transactions/$host_or_client.txt
    else 
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

      for (( i=host_or_client; i<=total; i=i+skip ))
      do
        ./app -exp=$exp -client=$i -config=configs/prod/node_$host_or_client.json < assets/data/transactions/$i.txt &
      done
    fi
  fi
fi
