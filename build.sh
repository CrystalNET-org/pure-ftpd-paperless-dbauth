#!/bin/bash
# Set the name of your Go script
SCRIPT_NAME="verify_pw.go"

# Set the output binary name
BINARY_NAME="verify_pw"

# Set the platform and architecture
PLATFORM="linux"
ARCHITECTURES=("amd64" "386" "arm64" "arm")

# Enable Go Modules
export GO111MODULE=on

mkdir ./out

for ARCH in "${ARCHITECTURES[@]}"
do
    env GOOS=${PLATFORM} GOARCH="${ARCH}" go mod tidy
    # Build the binary
    env GOOS=${PLATFORM} GOARCH="${ARCH}" go build -o "./out/${BINARY_NAME}_${ARCH}" ${SCRIPT_NAME}
    # Check if the build was successful
    if [ $? -eq 0 ]; then
        echo "Build successful! The binary is named '${BINARY_NAME}_${ARCH}' in the ./out directory"
    else
        echo "Build failed!"
    fi
done


