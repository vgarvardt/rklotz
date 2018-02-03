#!/bin/bash

set -e

source source-code/ci/scripts/helpers/prepare-gopath.sh

cd $SOURCE_GOPATH

echo "===> Run tests..."
make deps
make test
