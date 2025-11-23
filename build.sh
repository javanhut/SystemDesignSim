#!/bin/bash

echo "System Design Simulator - Build Script"
echo "======================================"
echo ""

echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "Go version: $(go version)"
echo ""

echo "Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "Error: Failed to download dependencies"
    exit 1
fi

echo ""
echo "Building application..."
echo "Note: First build may take 5-10 minutes due to Fyne compilation with CGo"
echo ""

go build -o systemdesignsim cmd/simulator/main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Build successful!"
    echo ""
    echo "Run the game with: ./systemdesignsim"
    echo ""
else
    echo ""
    echo "✗ Build failed"
    echo ""
    echo "Common issues:"
    echo "- Missing C compiler (gcc)"
    echo "- Missing graphics libraries (see docs/GETTING_STARTED.md)"
    echo "- Network issues downloading dependencies"
    exit 1
fi
