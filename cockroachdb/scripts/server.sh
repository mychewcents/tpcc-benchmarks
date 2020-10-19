#! /bin/sh

h=$(hostname -s | cut -c5-6)
action=$1

if [ $1 == "help" ]
then
  echo
  echo
  echo "This is simple manual for the executable."
  echo
  echo
  echo "\"cdbserv\" accepts two different types of arguments:"
  echo "   - start : To start the CockroachDB server on the localhost."
  echo "   - stop  : To stop the local running instance of the Cockroach DB."
  echo
  echo
  echo "Happy Running!"
elif [ $1 == "start" ]
then
  printf "**********\nStarting Cockroach DB node on : node${h}\n**********\n\n"
  if [ $h == 30 ]
  then
    cockroach start --insecure --store=/home/stuproj/cs4224m/crdb-node-files/node0 \
      --listen-addr=localhost:4050 \
      --http-addr=localhost:8080 \
      --join=xcnc30.comp.nus.edu.sg:4050,xcnc31.comp.nus.edu.sg:4050,xcnc32.comp.nus.edu.sg:4050,xcnc33.comp.nus.edu.sg:4050,xcnc34.comp.nus.edu.sg:4050 \
      --background

    # cockroach init --insecure --host=localhost:4050

  elif [ $h == 31 ]
  then
    cockroach start --insecure --store=/home/stuproj/cs4224m/crdb-node-files/node1 \
      --listen-addr=localhost:4050 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:4050,xcnc31.comp.nus.edu.sg:4050,xcnc32.comp.nus.edu.sg:4050,xcnc33.comp.nus.edu.sg:4050,xcnc34.comp.nus.edu.sg:4050 \
      --background

  elif [ $h == 32 ]
  then
    cockroach start --insecure --store=/home/stuproj/cs4224m/crdb-node-files/node2 \
      --listen-addr=localhost:4050 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:4050,xcnc31.comp.nus.edu.sg:4050,xcnc32.comp.nus.edu.sg:4050,xcnc33.comp.nus.edu.sg:4050,xcnc34.comp.nus.edu.sg:4050 \
      --background

  elif [ $h == 33 ]
  then
    cockroach start --insecure --store=/home/stuproj/cs4224m/crdb-node-files/node3 \
      --listen-addr=localhost:4050 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:4050,xcnc31.comp.nus.edu.sg:4050,xcnc32.comp.nus.edu.sg:4050,xcnc33.comp.nus.edu.sg:4050,xcnc34.comp.nus.edu.sg:4050 \
      --background

  elif [ $h == 34 ]
  then
    cockroach start --insecure --store=/home/stuproj/cs4224m/crdb-node-files/node4 \
      --listen-addr=localhost:4050 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:4050,xcnc31.comp.nus.edu.sg:4050,xcnc32.comp.nus.edu.sg:4050,xcnc33.comp.nus.edu.sg:4050,xcnc34.comp.nus.edu.sg:4050 \
      --background

  fi

  printf "\n**********\nStarted server on : node${h}\n**********\n"
elif [ $1 == "stop" ]
then
  printf "**********\nStopping the node : node${h}\n**********\n\n"
  cockroach quit --insecure --host=localhost:4050
  printf "\n**********\nStopped the node : node${h}\n**********\n"
elif [ $1 == "init" ]
then
  printf "**********\nInitializing the cluster : node${h}\n**********\n\n"
  cockroach init --insecure --host=localhost:4050
  printf "\n**********\nStopped the node : node${h}\n**********\n"
fi
