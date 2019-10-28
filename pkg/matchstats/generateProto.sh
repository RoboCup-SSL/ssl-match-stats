#!/bin/sh

# Fail on errors
set -e

# Update to latest protobuf compiler
go get -u github.com/golang/protobuf/protoc-gen-go

# Set package name to current directory
packageName=${PWD##*/}

# compile profobuf files in current directory
protoc -I. \
  -I${GOPATH}/src \
  -I${GOPATH}/src/github.com/RoboCup-SSL/ssl-go-tools/pkg/sslproto \
  --go_out=import_path="${packageName}:." \
  ./*.proto
