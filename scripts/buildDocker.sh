#!/usr/bin/env sh

set -e

source "$(dirname $0)/_variables.sh"

docker run --rm -v ${PWD}:/opt -it golang:1.20-alpine3.17 sh -c \
  "cd /opt && go build -o ${BINARY}"
