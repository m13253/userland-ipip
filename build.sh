#!/bin/bash

set -e
export GOOS=linux

mkdir -p build
go get -d -u ./cmd/ipip
go build -o build/ipip ./cmd/ipip

echo 'Built successfully to ./build/ipip'
