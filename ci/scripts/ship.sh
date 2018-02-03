#!/bin/bash

set -e

source source-code/ci/scripts/helpers/prepare-gopath.sh

cd $SOURCE_GOPATH

echo "===> Run goreleaser"
make deps
go get github.com/goreleaser/goreleaser
goreleaser --rm-dist
