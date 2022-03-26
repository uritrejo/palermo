#!/bin/bash

set -e

export GO111MODULE=on

#go vet ./...

go build cmd/palermo-server.go

go test ./...


exit 0
