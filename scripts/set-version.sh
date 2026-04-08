#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:?Usage: set-version.sh <version>}"

sed -i "s/Version: \"[^\"]*\"/Version: \"$VERSION\"/" cmd/maven-deps/main.go
sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION\"/" .claude-plugin/plugin.json

echo "Version set to $VERSION"
