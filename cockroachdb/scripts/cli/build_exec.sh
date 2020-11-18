#!/bin/bash

go build -o serverCmd cmd/server/main.go
echo "Example: ./serverCmd -env=prod -config=configs/prod/setup.json -node=1 (download-dataset | setup-dirs | start | stop | init | load | load-csv)"
echo "Example: ./serverCmd -env=prod -exp=5 -config=configs/prod/setup.json -node=1 run-exp"

go build -o clientCmd cmd/client/main.go
go build -o dbstateCmd cmd/dbstate/main.go
echo "Example: ./dbstateCmd -env=prod -config=configs/prod/setup.json -node=1"

chmod a+x scripts/export_data.sh
chmod a+x scripts/server.sh
chmod a+x scripts/run.sh
