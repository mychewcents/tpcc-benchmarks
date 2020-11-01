#! /bin/sh

exp=$1
server=$2
skip=5

if [ "$#" -eq 2 ] 
then
  if [ $exp -gt 4 ] && [ $exp -lt 9 ] 
  then

    host=192.168.49.179
    case $server in
      2) host=192.168.49.180
      ;;
      3) host=192.168.49.181
      ;;
      4) host=192.168.49.182
      ;;
      5) host=192.168.49.183
      ;;
    esac

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
      # echo $host $i
      ./app -exp=$exp -client=$i -host=$host -port=26257 < $i.txt &
    done
  fi
fi
