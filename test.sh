#!/bin/bash
set -e

# Initialize coverage file
echo "" > coverage.txt

# Ensure modules are up to date
go mod tidy

# Run tests with race detection and coverage
for pkg in $(go list ./...); do
    go test -race -coverprofile=profile.out -covermode=atomic $pkg
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done