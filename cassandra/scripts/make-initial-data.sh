#!/bin/bash
cd "$(dirname "$0")/.."

rm -rf assets/data

mkdir assets/data
mkdir assets/data/raw
mkdir assets/data/processed
mkdir assets/data/transactions

curl 'http://www.comp.nus.edu.sg/~cs4224/project-files.zip' -L -o assets/project-files.zip
unzip assets/project-files.zip -d assets
mv assets/project-files/data-files/* assets/data/raw
mv assets/project-files/xact-files/* assets/data/transactions

rm assets/project-files.zip
rm -rf assets/project-files

python3 scripts/python/make-initial-data.py
