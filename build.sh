#!/usr/bin/env bash

echo "Compiling application"
CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s"

echo "Build Docker image"
docker build -t jonasfred/sb-web .