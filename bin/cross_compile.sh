#!/bin/bash

script="$0"
FOLDER="$(pwd)/$(dirname $script)"

source $FOLDER/utils.sh
PROJECT_ROOT="$(abspath $FOLDER/..)"
echo "project root folder $PROJECT_ROOT"

echo "build docker image"
/bin/bash $FOLDER/build_golang.sh

##### RUN #####
echo "Starting container..."

# compile x86 binaries
for GOOS in darwin linux windows freebsd; do
  for GOARCH in 386 amd64; do
    docker run --rm \
               -it \
               -v $PROJECT_ROOT/go:/usr/src/slowport-restarter \
               -e GOOS=$GOOS \
               -e GOARCH=$GOARCH \
               -w /usr/src/slowport-restarter \
               golang:1.8 \
               bash -c "go build -v -o build/slowport-restarter-$GOOS-$GOARCH"
  done
done

# compile ARM binaries
for GOOS in linux darwin freebsd; do
  for GOARCH in arm; do
    docker run --rm \
               -it \
               -v $PROJECT_ROOT/go:/usr/src/slowport-restarter \
               -e GOOS=$GOOS \
               -e GOARCH=$GOARCH \
               -e GOARM=6 \
               -w /usr/src/slowport-restarter \
               golang:1.8 \
               bash -c "go build -v -o build/slowport-restarter-$GOOS-$GOARCH"
  done
done
