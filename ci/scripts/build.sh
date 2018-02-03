#!/bin/bash

set -e

source source-code/ci/scripts/helpers/prepare-gopath.sh

cd $SOURCE_GOPATH

echo "===> Dry-run goreleaser"
make deps
go get github.com/goreleaser/goreleaser
goreleaser --snapshot --rm-dist
