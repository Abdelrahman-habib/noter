#!/bin/bash
set -e

echo "Building application..."
go build -o bin/noter ./cmd/web

echo "Starting application..."
exec ./bin/noter "$@"
