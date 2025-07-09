#!/bin/bash

# Create Sentry Release Script
# This script creates a Sentry release and associates commits with it

set -e

# Check required environment variables
if [ -z "$SENTRY_AUTH_TOKEN" ]; then
    echo "Error: SENTRY_AUTH_TOKEN environment variable is required"
    exit 1
fi

if [ -z "$SENTRY_ORG" ]; then
    echo "Error: SENTRY_ORG environment variable is required"
    exit 1
fi

# Default values
SENTRY_PROJECT=${SENTRY_PROJECT:-"rune"}
VERSION=${1:-$(git describe --tags --always)}
ENVIRONMENT=${ENVIRONMENT:-"production"}

echo "Creating Sentry release: $VERSION"

# Create the release
curl -X POST \
  "https://sentry.io/api/0/organizations/$SENTRY_ORG/releases/" \
  -H "Authorization: Bearer $SENTRY_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"version\": \"$VERSION\",
    \"projects\": [\"$SENTRY_PROJECT\"],
    \"refs\": [{
      \"repository\": \"ferg-cod3s/rune\",
      \"commit\": \"$(git rev-parse HEAD)\"
    }]
  }"

echo "Sentry release $VERSION created successfully"

# Associate commits with the release
echo "Associating commits with release..."
curl -X POST \
  "https://sentry.io/api/0/organizations/$SENTRY_ORG/releases/$VERSION/commitfiles/" \
  -H "Authorization: Bearer $SENTRY_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"refs\": [{
      \"repository\": \"ferg-cod3s/rune\",
      \"commit\": \"$(git rev-parse HEAD)\"
    }]
  }"

echo "Commits associated with release $VERSION"

# Deploy the release to environment
echo "Deploying release to $ENVIRONMENT..."
curl -X POST \
  "https://sentry.io/api/0/organizations/$SENTRY_ORG/releases/$VERSION/deploys/" \
  -H "Authorization: Bearer $SENTRY_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"environment\": \"$ENVIRONMENT\",
    \"name\": \"$VERSION deployment\"
  }"

echo "Release $VERSION deployed to $ENVIRONMENT"