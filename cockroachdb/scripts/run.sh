#! /bin/sh

exp=$1
server=$2
skip=5

if [ "$#" -eq 2 ] 
then
  if [ $exp -gt 4 ] && [ $exp -lt 9 ] 
  then
    go build -o app ../cmd/app/main.go
    chmod a+x ./app
    case $exp in
      5) 
      clients=5
      skip=4
      ;;
      6) clients=4
      skip=5
      ;;
      7) clients=10
      skip=4
      ;;
      8) clients=8
      skip=5
      ;;
    esac

    tx=$server

    for (( i=tx; i<=clients*skip; i=i+skip ))
    do
      ./app -exp=$exp -client=$i < $i.txt &
    done
  fi
fi
