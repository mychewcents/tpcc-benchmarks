#!/bin/bash


go build -o setupCmd cmd/setup/main.go
echo "Example: ./setupCmd -env=prod -config=configs/prod/setup.json -node=1"

go build -o serverCmd cmd/server/*
echo "Example: ./serverCmd -env=prod -config=configs/prod/setup.json -node=1 (start | stop | init | load)"
echo "Example: ./serverCmd -env=prod -exp=5 -config=configs/prod/setup_5.json -node=1 run-exp"

go build -o clientCmd cmd/app/main.go
go build -o dbstateCmd cmd/dbstate/main.go
echo "Example: ./dbstateCmd -env=prod -config=configs/prod/setup.json -node=1"

chmod a+x scripts/init_setup.sh
chmod a+x scripts/server.sh
chmod a+x scripts/run.sh
