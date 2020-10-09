#! /bin/sh

arg=$1

if [ $1 == "all" ]
then
  echo "Stopping the entire cluster..."

  cockroach quit --host=xcnc30.comp.nus.edu.sg:30000
  cockroach quit --host=xcnc31.comp.nus.edu.sg:30001
  cockroach quit --host=xcnc32.comp.nus.edu.sg:30002
  cockroach quit --host=xcnc33.comp.nus.edu.sg:30003
  cockroach quit --host=xcnc34.comp.nus.edu.sg:30004
  
elif [ $1 == 0 ]
then
  echo "Stopping the node0..."

  cockroach quit --host=xcnc30.comp.nus.edu.sg:30000

elif [ $1 == 1 ]
then 
  echo "Stopping the node1..."

  cockroach quit --host=xcnc31.comp.nus.edu.sg:30001

elif [ $1 == 2 ]
then
  echo "Stopping the node2..."

  cockroach quit --host=xcnc32.comp.nus.edu.sg:30002

elif [ $1 == 3 ]
then
  echo "Stopping the node3..."

  cockroach quit --host=xcnc33.comp.nus.edu.sg:30003

elif [ $1 == 4 ]
then 
  echo "Stopping the node4..."

  cockroach quit --host=xcnc34.comp.nus.edu.sg:30004

fi

echo "Completed."
