#!/usr/bin/env sh

set -e

cd $(dirname $0)

if [ ! -f "$(basename $0)" ]; then
  echo "Incorrect location"
  exit 1
fi

source ./_variables.sh
cd ..

docker run --rm -v ${PWD}:/opt -it golang:1.20-alpine3.17 sh -c \
  "cd /opt && go build -o ${BINARY}"
