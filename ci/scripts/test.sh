#!/bin/sh

# Ensure this script fails if anything errors
set -e

set -ex

source source-code/ci/scripts/helpers/prepare-gopath.sh

cd $SOURCE_GOPATH

make deps
make test
