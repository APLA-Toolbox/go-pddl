#!/usr/bin/env bash

set -e

go mod vendor
go generate ./src

echo TEST : "${TEST}"
echo DEBUG : "${DEBUG}"
source .env

air
