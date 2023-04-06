#!/usr/bin/env bash

# checks that upx is installed we need to compress our binary to
# make it smaller, as router has limited storage
if [ -z "$(which upx)" ]; then
    echo "upx is not installed"
    exit 1
fi

# checks that go is installed, as nothing gonna work without it
if [ -z "$(which go)" ]; then
    echo "go is not installed"
    exit 1
fi

# define list of architectures to build for. add your own
# architecture there if it's not covered by this list
ARCHITECTURES="mips mipsle mips64 mips64le amd64 arm arm64 "

# build for each architecture
for ARCH in $ARCHITECTURES; do
    echo "Building for $ARCH"
    # cross-compile for $ARCH
    GOOS=linux GOARCH=$ARCH go build -o bin/network-status-$ARCH

    # minimise the binary
    upx -9 bin/network-status-$ARCH
done
