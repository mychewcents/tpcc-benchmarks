#!/bin/bash


go build -o setupCmd cmd/setup/main.go
echo "Example: ./setupCmd -env=prod -node=1 -config=configs/prod/setup_5.json"

go build -o serverCmd cmd/server/main.go
echo "Example: ./serverCmd -env=prod -node=1 -config=configs/prod/setup_5.json (start | stop | init | run-exp | load)"

go build -o clientCmd cmd/app/main.go
go build -o dbstateCmd cmd/dbstate/main.go
echo "Example: ./dbstateCmd -env=prod -node=1 -config=configs/prod/setup_5.json"

chmod a+x scripts/init_setup.sh
chmod a+x scripts/server.sh
chmod a+x scripts/run.sh