#!/bin/bash

export CGO_ENABLED=0 
export GOPATH=$(pwd)
export GOOS=linux
export GOARCH=amd64
go build -a -installsuffix cgo
docker build -t bestbug/weatherai .