#!/bin/bash


go build -o setup cmd/setup/main.go
echo "./setup -env=prod -node=1 -config=configs/prod/setup_5.json"

go build -o server cmd/server/main.go
echo "./server -env=prod -node=1 -config=configs/prod/setup_5.json (start | stop | init | sql | load)"

go build -o init cmd/init/main.go
go build -o client cmd/app/main.go
go build -o dbstate cmd/dbstate/main.go