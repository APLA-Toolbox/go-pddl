#!/usr/bin/env bash

if [[ "$TEST" -eq "1" ]] && [[ "$DEBUG" -eq "1" ]]; then
    echo "Debugging tests..."
    dlv test --listen=:2345 --headless=true --api-version=2 -p 1 -v ./...
elif [[ "$TEST" -eq "1" ]]; then
    echo "Starting all tests..."
    go test -p 1 -v ./...
elif [[ "$DEBUG" -eq 1 ]]; then
    echo "Debugging..."
    dlv --listen=:2345 --headless=true --api-version=2 exec ./tmp/main
else
    echo "Starting..."
    ./tmp/main
fi
