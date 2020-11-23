#!/usr/bin/env bash

set -e

export BUILD_TIME=$(date '+%FT%T%:z')
export BUILD_GIT=$(git rev-parse HEAD)

go mod vendor
go generate ./src

echo TEST : "${TEST}"
echo DEBUG : "${DEBUG}"
echo TEST_DEBUG_PACKAGE : "${TEST_DEBUG_PACKAGE}"

source .env

air 