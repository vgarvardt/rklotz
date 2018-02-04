#!/bin/bash

set -e

source source-code/ci/scripts/helpers/prepare-gopath.sh

cd $SOURCE_GOPATH

echo "===> Run code coverage"
make deps
go get github.com/go-playground/overalls
overalls -project=github.com/vgarvardt/rklotz -covermode=count
bash <(curl -s https://codecov.io/bash) -f overalls.coverprofile
