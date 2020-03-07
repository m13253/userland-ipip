#!/bin/bash

set -e
export GOOS=linux

mkdir -p build
go build -o build/ipip ./cmd/ipip

echo 'Built successfully to ./build/ipip'
