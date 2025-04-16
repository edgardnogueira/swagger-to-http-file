#!/bin/bash

# Script to install Git hooks for swagger-to-http-file
#
# This script:
# 1. Creates a directory for the hooks if it doesn't exist
# 2. Copies the hook scripts from the scripts/hooks directory to the .git/hooks directory
# 3. Makes the hooks executable

set -e

# Determine the script directory and the Git root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GIT_ROOT="$(git rev-parse --show-toplevel)"
HOOKS_DIR="$GIT_ROOT/.git/hooks"
SOURCE_HOOKS_DIR="$SCRIPT_DIR/hooks"

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Install hooks
echo "Installing Git hooks..."
for hook in "$SOURCE_HOOKS_DIR"/*; do
  if [ -f "$hook" ]; then
    hook_name=$(basename "$hook")
    target="$HOOKS_DIR/$hook_name"
    cp "$hook" "$target"
    chmod +x "$target"
    echo "Installed $hook_name hook"
  fi
done

echo "Git hooks installation complete!"
echo "The following hooks are now active:"
ls -la "$HOOKS_DIR"

echo ""
echo "You can configure the hooks by setting the following environment variables:"
echo "  SWAGGER_TO_HTTP_SKIP_HOOKS=1    # Skip running the Git hooks"
echo "  SWAGGER_TO_HTTP_FILES=\"file1.json file2.yaml\"    # Specific Swagger/OpenAPI files to check"
echo "  SWAGGER_TO_HTTP_OUTPUT_DIR=\"./http\"    # Directory for generated HTTP files"
echo ""
