#! /bin/sh

#############################################
# The init.sh script initializes the cdbserv 
# command line tool and downloads the 
# project files for data and transaction.
#############################################

if [ ! -d "/temp/cs5424-team-m/cdb-server" ]
then
  echo "Creating the server installation script directory"
  mkdir /temp/cs5424-team-m/cdb-server
fi

cp scripts/server.sh scripts/cdbserv
chmod a+x scripts/cdbserv
mv scripts/cdbserv /temp/cs5424-team-m/cdb-server/

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
