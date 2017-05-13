#!/bin/bash

script="$0"
FOLDER="$(pwd)/$(dirname $script)"

source $FOLDER/utils.sh
PROJECT_ROOT="$(abspath $FOLDER/..)"
echo "project root folder $PROJECT_ROOT"

echo "build docker image"
/bin/bash $FOLDER/build_golang.sh

##### VOLUMES #####

##### RUN #####
echo "Starting container..."

GOOS=darwin
GOARCH=amd64

docker run --rm \
           -it \
           -v $PROJECT_ROOT/go:/usr/src/slowport-restarter \
           -e GOOS=$GOOS \
           -e GOARCH=$GOARCH \
           -w /usr/src/slowport-restarter \
           golang:1.8 \
           bash -c "go build -v -o build/slowport-restarter-$GOOS-$GOARCH"
