#!/bin/bash

set -e

export GO111MODULE=on

#go vet ./...

go build cmd/main.go

go test ./...


exit 0
