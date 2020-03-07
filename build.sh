#!/bin/bash

export GOOS=linux

mkdir -p build
go build -o build/ipip ./cmd/ipip
