#!/bin/bash

# Build script that ensures telemetry keys are embedded
# This script can be used for local builds and CI/CD

set -e

# Default telemetry keys (can be overridden by environment variables)
DEFAULT_SEGMENT_KEY="ZkEZXHRWH96y8EviNkbYJUByqGR9QI4G"
DEFAULT_SENTRY_DSN="https://3b20acb23bbbc5958448bb41900cdca2@sentry.fergify.work/10"

# Use environment variables if set, otherwise use defaults
SEGMENT_KEY=${RUNE_SEGMENT_WRITE_KEY:-$DEFAULT_SEGMENT_KEY}
SENTRY_DSN=${RUNE_SENTRY_DSN:-$DEFAULT_SENTRY_DSN}
VERSION=${VERSION:-$(git describe --tags --always --dirty)}

echo "Building Rune with telemetry support..."
echo "Version: $VERSION"
echo "Segment Key: ${SEGMENT_KEY:0:10}..."
echo "Sentry DSN: ${SENTRY_DSN:0:30}..."

# Build the binary
go build -ldflags "\
  -s -w \
  -X github.com/ferg-cod3s/rune/internal/commands.version=$VERSION \
  -X github.com/ferg-cod3s/rune/internal/telemetry.version=$VERSION \
  -X github.com/ferg-cod3s/rune/internal/telemetry.segmentWriteKey=$SEGMENT_KEY \
  -X github.com/ferg-cod3s/rune/internal/telemetry.sentryDSN=$SENTRY_DSN \
" -o rune ./cmd/rune

echo "Build completed successfully!"
echo "Binary: ./rune"

# Test the build
echo "Testing telemetry integration..."
if RUNE_DEBUG=true ./rune --version > /dev/null 2>&1; then
    echo "✅ Telemetry integration test passed"
else
    echo "❌ Telemetry integration test failed"
    exit 1
fi