#!/bin/bash
set -e

go mod tidy

golangci-lint run --fix

go test ./...
