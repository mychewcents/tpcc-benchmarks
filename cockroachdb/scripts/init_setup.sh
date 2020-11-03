#!/bin/bash

############################################################
# The init.sh script initializes the cdbserv and cdbclient
# command line tools and downloads the project files for data 
# and transaction.
############################################################

if [[ $1 == "help" ]]
then
  echo
  echo
  echo "This is a simple manual for the \"init_setup.sh\" file."
  echo
  echo
  echo "\"init_script.sh\" accepts 1 command line argument"
  echo
  echo
  echo "First Argument - Type of the environment"
  echo "   - prod   : To install the \"cdbserv\" and \"cdbclient\" exec "
  echo "              at the default location and download the project files"
  echo "   - local  : Only downloads the project files"
  echo
  echo
  echo "Second Argument - Location of the working directory"
  echo "   - <path>   : Path to the working directory for the nodes and executables"
  echo
  echo
  echo "Third Argument - URL for the project file downloads"
  echo "   - <URL>    : To download the files for the project"
  echo
  echo
  echo "Happy Running!"
  echo
  echo
elif [[ "$#" -eq 4 ]]
then
  if [[ $1 == 'prod' ]] || [[ $1 == 'local' ]]
  then
    if [[ ! -d "${2}/cdb-server" ]]
      then
        echo "Creating the server installation script directory"
        mkdir -p $2/cdb-server
      fi

    extern_dir=$2/cdb-server/node-files/$4/extern 
    rm -rf $2/cdb-server/node-files/$4
    mkdir -p $extern_dir/assets/raw
    mkdir -p $extern_dir/assets/processed
    cp assets/data/raw/* $extern_dir/assets/raw
    cp -r assets/data/processed/* $extern_dir/assets/processed
  else
    echo "Use the \"help\" command to learn more about the arguments"
  fi
else
    echo "Use the \"help\" command to learn more about the arguments"
fi