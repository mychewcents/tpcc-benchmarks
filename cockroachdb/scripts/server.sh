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
    cockroach start --insecure --store=/temp/cs5424-team-m/node0 \
      --listen-addr=localhost:26257 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:26257,xcnc31.comp.nus.edu.sg:26257,xcnc32.comp.nus.edu.sg:26257,xcnc33.comp.nus.edu.sg:26257,xcnc34.comp.nus.edu.sg:26257 \
      --background

    # cockroach init --insecure --host=localhost:30000

  elif [ $h == 31 ]
  then
    cockroach start --insecure --store=/temp/cs5424-team-m/node1 \
      --listen-addr=localhost:26257 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:26257,xcnc31.comp.nus.edu.sg:26257,xcnc32.comp.nus.edu.sg:26257,xcnc33.comp.nus.edu.sg:26257,xcnc34.comp.nus.edu.sg:26257 \
      --background

  elif [ $h == 32 ]
  then
    cockroach start --insecure --store=/temp/cs5424-team-m/node2 \
      --listen-addr=localhost:30000 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc31.comp.nus.edu.sg:30000,xcnc32.comp.nus.edu.sg:30000,xcnc33.comp.nus.edu.sg:30000,xcnc34.comp.nus.edu.sg:30000 \
      --background

  elif [ $h == 33 ]
  then
    cockroach start --insecure --store=/temp/cs5424-team-m/node3 \
      --listen-addr=localhost:30000 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc31.comp.nus.edu.sg:30000,xcnc32.comp.nus.edu.sg:30000,xcnc33.comp.nus.edu.sg:30000,xcnc34.comp.nus.edu.sg:30000 \
      --background

  elif [ $h == 34 ]
  then
    cockroach start --insecure --store=/temp/cs5424-team-m/node4 \
      --listen-addr=localhost:30000 \
      --http-addr=localhost:40000 \
      --join=xcnc30.comp.nus.edu.sg:30000,xcnc31.comp.nus.edu.sg:30000,xcnc32.comp.nus.edu.sg:30000,xcnc33.comp.nus.edu.sg:30000,xcnc34.comp.nus.edu.sg:30000 \
      --background

  fi

  printf "\n**********\nStarted server on : node${h}\n**********\n"
elif [ $1 == "stop" ]
then
  printf "**********\nStopping the node : node${h}\n**********\n\n"
  cockroach quit --insecure --host=localhost:26257
  printf "\n**********\nStopped the node : node${h}\n**********\n"
elif [ $1 == "init" ]
then
  printf "**********\nInitializing the cluster : node${h}\n**********\n\n"
  cockroach init --insecure --host=localhost:26257
  printf "\n**********\nStopped the node : node${h}\n**********\n"
fi
