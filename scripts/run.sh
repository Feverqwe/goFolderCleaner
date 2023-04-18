#!/usr/bin/env sh

set -e

cd $(dirname $0)

if [ ! -f "$(basename $0)" ]; then
  echo "Incorrect location"
  exit 1
fi

if [ "$1" = "dev" ]; then
    export DEBUG_UI=1
fi

source ./_variables.sh
cd ..

sh ./scripts/build.sh
./$BINARY
