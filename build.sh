#!/usr/bin/env bash
set -xe

# install packages and dependencies
go get github.com/gorilla/handlers
go get github.com/gorilla/mux
go get github.com/aws/aws-sdk-go/...

go mod init fam-photos-server
go install fam-photos-server

mkdir -p bin

# build command
GOARCH=amd64 GOOS=linux go build -o bin/application application.go
