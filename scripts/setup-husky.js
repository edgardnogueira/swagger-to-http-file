#!/usr/bin/env node

/**
 * Setup script for integrating Swagger-to-HTTP with Husky
 * 
 * This script:
 * 1. Checks if Husky is installed in the project
 * 2. If not, suggests installation commands
 * 3. If it is, creates Husky hook scripts for pre-commit and post-checkout
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Determine if we're in a git repository
let gitRoot;
try {
  gitRoot = execSync('git rev-parse --show-toplevel').toString().trim();
} catch (error) {
  console.error('Error: Not in a Git repository.');
  process.exit(1);
}

// Check if Husky is installed
const huskyDir = path.join(gitRoot, '.husky');
const nodeModulesHusky = path.join(gitRoot, 'node_modules', 'husky');
const packageJsonPath = path.join(gitRoot, 'package.json');

let hasHusky = fs.existsSync(huskyDir);
let hasHuskyNodeModules = fs.existsSync(nodeModulesHusky);
let hasPackageJson = fs.existsSync(packageJsonPath);

if (!hasHusky || !hasHuskyNodeModules) {
  console.log('Husky doesn\'t appear to be installed or configured in this project.');
  console.log('\nTo install Husky:');
  
  if (!hasPackageJson) {
    console.log('1. Initialize package.json first:');
    console.log('   npm init -y');
  }
  
  console.log('2. Install Husky:');
  console.log('   npm install husky --save-dev');
  console.log('   npx husky install');
  console.log('   npm pkg set scripts.prepare="husky install"');
  console.log('\nAfter installing Husky, run this script again.');
  process.exit(0);
}

// Create Husky hook scripts
console.log('Setting up Husky hooks for Swagger-to-HTTP...');

// Pre-commit hook
const preCommitPath = path.join(huskyDir, 'pre-commit');
const preCommitContent = `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Skip if SWAGGER_TO_HTTP_SKIP_HOOKS is set
if [ -n "$SWAGGER_TO_HTTP_SKIP_HOOKS" ]; then
  echo "Skipping Swagger-to-HTTP hooks (SWAGGER_TO_HTTP_SKIP_HOOKS is set)"
  exit 0
fi

# Check if swagger-to-http-file is installed
if ! command -v swagger-to-http-file &> /dev/null; then
  echo "Warning: swagger-to-http-file not found in PATH. Skipping HTTP file generation."
  echo "To install: go install github.com/edgardnogueira/swagger-to-http-file/cmd/swagger-to-http-file@latest"
  exit 0
fi

# Set default values
OUTPUT_DIR="\${SWAGGER_TO_HTTP_OUTPUT_DIR:-.}"

# Get staged Swagger files
SWAGGER_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '\\.(json|yaml|yml)$' || true)

if [ -n "$SWAGGER_FILES" ]; then
  echo "Processing Swagger files: $SWAGGER_FILES"
  
  for file in $SWAGGER_FILES; do
    echo "Converting $file..."
    swagger-to-http-file -i "$file" -o "$OUTPUT_DIR" -w
    
    # Add generated HTTP files to the commit
    for http_file in $(find "$OUTPUT_DIR" -name "*.http" -type f -newer "$file"); do
      echo "Adding generated file to commit: $http_file"
      git add "$http_file"
    done
  done
fi
`;

// Post-checkout hook
const postCheckoutPath = path.join(huskyDir, 'post-checkout');
const postCheckoutContent = `#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Skip if not a branch checkout or if skip flag is set
if [ "$3" != "0" ] || [ -n "$SWAGGER_TO_HTTP_SKIP_HOOKS" ]; then
  exit 0
fi

# Check if swagger-to-http-file is installed
if ! command -v swagger-to-http-file &> /dev/null; then
  echo "Warning: swagger-to-http-file not found in PATH. Skipping HTTP file generation."
  exit 0
fi

# Set default values
OUTPUT_DIR="\${SWAGGER_TO_HTTP_OUTPUT_DIR:-.}"

# Get changed Swagger files between branches
PREV_HEAD="$1"
NEW_HEAD="$2"

SWAGGER_FILES=$(git diff --name-only "$PREV_HEAD" "$NEW_HEAD" | grep -E '\\.(json|yaml|yml)$' || true)

if [ -n "$SWAGGER_FILES" ]; then
  echo "Processing changed Swagger files between branches:"
  echo "$SWAGGER_FILES"
  
  for file in $SWAGGER_FILES; do
    if [ -f "$file" ]; then
      echo "Converting $file..."
      swagger-to-http-file -i "$file" -o "$OUTPUT_DIR" -w
    fi
  done
fi
`;

// Write the hook files
try {
  fs.writeFileSync(preCommitPath, preCommitContent);
  fs.chmodSync(preCommitPath, '755');
  console.log('✓ Created pre-commit hook');
  
  fs.writeFileSync(postCheckoutPath, postCheckoutContent);
  fs.chmodSync(postCheckoutPath, '755');
  console.log('✓ Created post-checkout hook');
  
  console.log('\nHusky hooks for Swagger-to-HTTP installed successfully!');
  console.log('\nYou can configure the hooks by setting the following environment variables:');
  console.log('  SWAGGER_TO_HTTP_SKIP_HOOKS=1    # Skip running the Git hooks');
  console.log('  SWAGGER_TO_HTTP_OUTPUT_DIR="./http"    # Directory for generated HTTP files');
} catch (error) {
  console.error(`Error creating hook files: ${error.message}`);
  process.exit(1);
}
