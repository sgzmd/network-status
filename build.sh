#!/usr/bin/env bash

# checks that upx is installed
if [ -z "$(which upx)" ]; then
    echo "upx is not installed"
    exit 1
fi

# checks that go is installed
if [ -z "$(which go)" ]; then
    echo "go is not installed"
    exit 1
fi

# define list of architectures to build for
ARCHITECTURES="mips mipsle mips64 mips64le amd64 arm arm64 "

# build for each architecture
for ARCH in $ARCHITECTURES; do
    echo "Building for $ARCH"
    GOOS=linux GOARCH=$ARCH go build -o bin/network-status-$ARCH
    upx -9 bin/network-status-$ARCH
done
