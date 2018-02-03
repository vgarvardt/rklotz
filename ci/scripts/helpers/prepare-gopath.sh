#!/bin/bash

set -e

# Include dotfiles to bash globs
shopt -s dotglob

echo "===> Setup GOPATH..."
export GOPATH=$PWD/gopath
export SOURCE_GOPATH=$GOPATH/src/github.com/$OWNER/$REPO
export PATH=${GOPATH}/bin:$PATH
mkdir -p ${GOPATH}/bin

echo "===> Move sources to GOPATH..."
[ -d $SOURCE_GOPATH ] && rm -rf $SOURCE_GOPATH
mkdir -p $SOURCE_GOPATH/
mv source-code/* $SOURCE_GOPATH/
