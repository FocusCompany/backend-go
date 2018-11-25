#!/usr/bin/env bash

set -e

# https://github.com/docker-library/postgres/pull/440
export DB_ADDR=/var/run/postgresql

PATH="$PATH:$GOPATH/bin/"
make migrate