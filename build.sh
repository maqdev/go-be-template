#!/bin/bash
set -e
set -x

# Check if the script was provided with exactly one argument
if [ $# -ne 1 ]; then
  echo "Usage: $0 path to version package e.g. github.com/maqdev/go-be-template/config"
  exit 1
fi

# Store the argument in a variable
version_pkg="$1"

# Put installed packages into ./bin
GOBIN=$PWD/`dirname $0`/bin
BRANCH="${STAMP_BRANCH}"
BUILD="${STAMP_VERSION}"

if [[ "${BUILD}" == "" && -d ".git" ]]; then
    BUILD=`git rev-parse --short HEAD || ""`
    BRANCH=`(git symbolic-ref --short HEAD | tr -d \/ ) || ""`
fi

if [[ "$BRANCH" = main || "$BRANCH" = master ]]; then
    BRANCH=""
fi

FLAGS="-X $version_pkg.branch=$BRANCH -X $version_pkg.build=$BUILD"

mkdir -p bin

# Build the binary passed in by arguments, defaulting to all cmds
go build -trimpath -ldflags "$FLAGS" -v -o "bin/" "${@:-./cmd/...}"
