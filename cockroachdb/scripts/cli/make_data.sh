#!/bin/bash

rm -rf assets/data

mkdir assets/data
mkdir assets/data/raw
mkdir assets/data/transactions
mkdir assets/data/processed
mkdir assets/data/processed/warehouse
mkdir assets/data/processed/district
mkdir assets/data/processed/customer
mkdir assets/data/processed/order
mkdir assets/data/processed/orderline
mkdir assets/data/processed/stock
mkdir assets/data/processed/item
mkdir assets/data/processed/itempairs

curl $1 -L -o assets/project-files.zip
unzip assets/project-files.zip -d assets
mv assets/project-files/data-files/* assets/data/raw
mv assets/project-files/xact-files/* assets/data/transactions

rm assets/project-files.zip
rm -rf assets/project-files