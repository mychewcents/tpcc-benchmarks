#! /bin/sh

#############################################
# The init.sh script initializes the cdbserv 
# command line tool and downloads the 
# project files for data and transaction.
#############################################

if [ "$#" -eq 0 ]
then
  if [ ! -d "/temp/cs5424-team-m/cdb-server" ]
  then
    echo "Creating the server installation script directory"
    mkdir /temp/cs5424-team-m/cdb-server
  fi

  cp scripts/server.sh cdbserv
  chmod a+x cdbserv
  mv cdbserv /temp/cs5424-team-m/cdb-server/
fi

rm -rf assets/data

mkdir assets/data
mkdir assets/data/raw
mkdir assets/data/transactions

curl 'http://www.comp.nus.edu.sg/~cs4224/project-files.zip' -L -o assets/project-files.zip
unzip assets/project-files.zip -d assets
mv assets/project-files/data-files/* assets/data/raw
mv assets/project-files/xact-files/* assets/data/transactions

rm assets/project-files.zip
rm -rf assets/project-files
