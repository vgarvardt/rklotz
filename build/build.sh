#!/bin/bash

set -e

# Get rid of existing binaries
rm -f dist/rklotz*

# Check if VERSION variable set and not empty, otherwise set to default value
if [ -z "$VERSION" ]; then
  VERSION="0.0.1-dev"
fi
echo "Building application version $VERSION"

echo "Building default binary"
CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X github.com/vgarvardt/rklotz/app.version=${VERSION}" -o "dist/rklotz" $PKG_SRC

# Build binaries
OS_PLATFORM_ARG=(linux darwin)
OS_ARCH_ARG=(amd64)
for OS in ${OS_PLATFORM_ARG[@]}; do
  for ARCH in ${OS_ARCH_ARG[@]}; do
    echo "Building binary for $OS/$ARCH..."
    GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X github.com/vgarvardt/rklotz/app.version=${VERSION}" -o "dist/rklotz.$OS.$ARCH" $PKG_SRC
  done
done
