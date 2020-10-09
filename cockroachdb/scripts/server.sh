#! /bin/sh

h=$(hostname -s | cut -c5-6)
action=$1

if [ $1 == "start" ]
then
  echo "Starting Cockroach DB node on : node${h}"
  if [ $h == 30 ]
  then
    cockroach start --insecure --store=node0 \
      --listen-addr=localhost:30000 \
      --http-addr=localhost:40000 \
      --join=xcnc31.comp.nus.edu.sg:30001,xcnc32.comp.nus.edu.sg:30002,xcnc33.comp.nus.edu.sg:30003,xcnc34.comp.nus.edu.sg:30004 \
      --background

    cockroach init --insecure --host=localhost:30000

  elif [ $h == 31 ]
  then
    cockroach start --insecure --store=node1 \
      --listen-addr=localhost:30001 \
      --http-addr=localhost:40001 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc32.comp.nus.edu.sg:30002,xcnc33.comp.nus.edu.sg:30003,xcnc34.comp.nus.edu.sg:30004 \
      --background

  elif [ $h == 32 ]
  then
    cockroach start --insecure --store=node2 \
      --listen-addr=localhost:30002 \
      --http-addr=localhost:40002 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc31.comp.nus.edu.sg:30001,xcnc33.comp.nus.edu.sg:30003,xcnc34.comp.nus.edu.sg:30004 \
      --background

  elif [ $h == 33 ]
  then
    cockroach start --insecure --store=node3 \
      --listen-addr=localhost:30003 \
      --http-addr=localhost:40003 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc31.comp.nus.edu.sg:30001,xcnc32.comp.nus.edu.sg:30002,xcnc34.comp.nus.edu.sg:30004 \
      --background

  elif [ $h == 34 ]
  then
    cockroach start --insecure --store=node4 \
      --listen-addr=localhost:30004 \
      --http-addr=localhost:40004 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc31.comp.nus.edu.sg:30001,xcnc32.comp.nus.edu.sg:30002,xcnc33.comp.nus.edu.sg:30003 \
      --background

  fi

  echo "\n\n*********\n\nStarted server on : node${h}"
elif [ $1 == "stop" ]
then
  echo "Stopping the node : node${h}..."
  if [ $h == 30 ]
  then
    cockroach quit --insecure --host=localhost:30000
  elif [ $h == 31 ]
  then
    cockroach quit --insecure --host=localhost:30001
  elif [ $h == 32 ]
  then
    cockroach quit --insecure --host=localhost:30002
  elif [ $h == 33 ]
  then
    cockroach quit --insecure --host=localhost:30003
  elif [ $h == 34 ]
  then
    cockroach quit --insecure --host=localhost:30004
  fi

fi
