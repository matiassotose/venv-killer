#!/bin/bash

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "Error: goreleaser is not installed."
    echo "Please install it using: brew install goreleaser/tap/goreleaser"
    exit 1
fi

# Check if GITHUB_TOKEN is set
if [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: GITHUB_TOKEN is not set."
    echo "Please export GITHUB_TOKEN with repo permissions."
    exit 1
fi

# Get the version to release
if [ -z "$1" ]; then
    echo "Usage: ./release.sh <version>"
    echo "Example: ./release.sh v0.1.0"
    exit 1
fi

VERSION=$1

# Confirm release
read -p "Are you sure you want to release version $VERSION? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Release cancelled."
    exit 1
fi

# Tag and push
echo "Creating tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

# Run goreleaser
echo "Running goreleaser..."
goreleaser release --clean

echo "Release $VERSION completed successfully!"
