#!/bin/bash


# Set the name of your Go script
SCRIPT_NAME="verify_pw.go"

# Set the output binary name
BINARY_NAME="verify_pw"

# Set the platform and architecture
PLATFORM="linux"
ARCHITECTURE="amd64"

# Enable Go Modules
export GO111MODULE=on

mkdir ./out

env GOOS=$PLATFORM GOARCH=$ARCHITECTURE go mod tidy

# Build the binary
env GOOS=$PLATFORM GOARCH=$ARCHITECTURE go build -o ./out/$BINARY_NAME $SCRIPT_NAME

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! The binary is named '$BINARY_NAME'"
else
    echo "Build failed!"
fi
