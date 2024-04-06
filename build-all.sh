#!/usr/bin/env bash

version=$1

if [ -z "$version" ]; then
  echo "Usage: build-all.sh <version>"
  exit
fi

GOOS=darwin GOARCH=amd64 go build
zip lace-$version-mac-amd64.zip lace
GOOS=linux GOARCH=amd64 go build
zip lace-$version-linux-amd64.zip lace
GOOS=windows GOARCH=amd64 go build
zip lace-$version-win-amd64.zip lace.exe
