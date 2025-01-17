#!/bin/bash

PACKAGE_NAME="endar_agent"
OUTPUT_DIR="build"

# Create output directory
mkdir -p $OUTPUT_DIR

# List of platforms to build for (OS/ARCH)
platforms=(
    "windows/amd64:_win64.exe"
    "windows/386:_win32.exe"
    "linux/amd64:_linux_amd64"
    "linux/arm64:_linux_arm64"
    "darwin/amd64:_macos_amd64"
    "darwin/arm64:_macos_arm64"
)

for platform in "${platforms[@]}"
do
    IFS=':' read -r os_arch suffix <<< "$platform"
    OS=${os_arch%/*}
    ARCH=${os_arch#*/}
    OUTPUT_NAME="$PACKAGE_NAME$suffix"

    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_DIR/$OUTPUT_NAME main.go

    if [ $? -ne 0 ]; then
        echo "Failed to build for $OS/$ARCH"
    else
        echo "Successfully built: $OUTPUT_NAME"
    fi
done

echo "All builds completed!"
