#!/bin/bash

set -ex

export GO111MODULE=on

#go vet ./...

suffix=$(go env GOEXE)
go build -o bin/palermo${suffix} cmd/main.go

go test ./...


exit 0
