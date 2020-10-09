#! /bin/sh

if [ ! -d "/temp/cs5424-team-m/cdb-server" ]
then
  echo "Creating the server installation script directory"
  mkdir /temp/cs5424-team-m/cdb-server
fi

cp server.sh cdbserv
chmod a+x cdbserv
mv cdbserv /temp/cs5424-team-m/cdb-server/
