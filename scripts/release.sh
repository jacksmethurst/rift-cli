#!/bin/bash

set -e

if [ -z "$1" ]; then
    echo "Usage: ./scripts/release.sh <version>"
    echo "Example: ./scripts/release.sh 1.0.6"
    exit 1
fi

VERSION=$1

# Validate version format
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format x.y.z (e.g., 1.0.6)"
    exit 1
fi

echo "🚀 Creating release $VERSION..."

# Update all version numbers
echo "📝 Updating version numbers..."
node scripts/update-version.js $VERSION

# Commit changes
echo "💾 Committing version updates..."
git add .
git commit -m "Bump version to $VERSION"

# Push changes
echo "📤 Pushing changes..."
git push origin main

# Create and push tag
echo "🏷️  Creating tag v$VERSION..."
git tag v$VERSION
git push origin v$VERSION

echo "✅ Release $VERSION created! Check GitHub Actions for build progress."
echo "🔗 https://github.com/jacksmethurst/rift-cli/actions"