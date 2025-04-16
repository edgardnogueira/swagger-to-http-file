# Git Hooks Integration

This document explains how to configure Git hooks for automatically updating HTTP files when Swagger/OpenAPI files are modified.

## Benefits of Using Git Hooks

- **Automatic updates**: HTTP files are always in sync with your Swagger/OpenAPI files
- **Simplified workflow**: No need to manually run the conversion tool
- **Prevent outdated files**: Ensures that HTTP files are updated when API specs change
- **Team collaboration**: Everyone on the team works with up-to-date HTTP files

## Quick Setup

### Option 1: Native Git Hooks

1. Run the installation script:
   ```bash
   # Make the script executable
   chmod +x scripts/install-hooks.sh
   
   # Run the installation script
   ./scripts/install-hooks.sh
   ```

2. The script will copy the hooks to your `.git/hooks` directory and make them executable.

### Option 2: Husky (for Node.js projects)

1. Run the Node.js setup script:
   ```bash
   # Make the script executable
   chmod +x scripts/setup-husky.js
   
   # Run the setup script
   node scripts/setup-husky.js
   ```

2. If Husky isn't installed, the script will provide instructions.

3. After installation, the script will create the necessary hook files.

## Configuration

You can configure the Git hooks using environment variables:

- `SWAGGER_TO_HTTP_SKIP_HOOKS=1`: Skip running the Git hooks
- `SWAGGER_TO_HTTP_FILES="file1.json file2.yaml"`: Specific Swagger/OpenAPI files to check
- `SWAGGER_TO_HTTP_OUTPUT_DIR="./http"`: Directory for generated HTTP files

## How It Works

### Pre-commit Hook

The pre-commit hook:
1. Checks if any Swagger/OpenAPI files are about to be committed
2. If so, runs the conversion tool to update the corresponding HTTP files
3. Adds the generated HTTP files to the commit

This ensures that whenever you commit a change to a Swagger file, the updated HTTP files are included in the same commit.

### Post-checkout Hook

The post-checkout hook:
1. Runs when you switch branches
2. Checks if any Swagger/OpenAPI files changed between branches
3. If so, updates the corresponding HTTP files

This ensures that when you switch branches, your HTTP files are updated to match the Swagger files in the new branch.

## Manual Synchronization

If you need to manually check for changes and update HTTP files:

```bash
./scripts/detect-swagger-changes.sh
```

This script:
1. Finds all Swagger/OpenAPI files in the current directory
2. Checks if they've been modified since their corresponding HTTP files
3. If so, runs the conversion tool to update the HTTP files

For more options, run:

```bash
./scripts/detect-swagger-changes.sh --help
```

## Integration with CI/CD

You can add the following step to your CI/CD pipeline to ensure HTTP files are up to date:

```yaml
- name: Check HTTP files are up to date
  run: |
    # Install swagger-to-http-file
    go install github.com/edgardnogueira/swagger-to-http-file/cmd/swagger-to-http-file@latest
    
    # Run the detection script
    ./scripts/detect-swagger-changes.sh --verbose
    
    # Check if any files were changed
    if git status --porcelain | grep -q ".http$"; then
      echo "ERROR: HTTP files are not up to date with Swagger files"
      git status --porcelain | grep ".http$"
      exit 1
    fi
```

## Integration with NPM Scripts

For Node.js projects, you can add the following NPM scripts to your `package.json`:

```json
{
  "scripts": {
    "swagger:convert": "swagger-to-http-file -i ./api/swagger.json -o ./http -w",
    "swagger:watch": "nodemon --watch ./api --ext json,yaml,yml --exec 'npm run swagger:convert'"
  }
}
```

Then you can run:

```bash
# Convert Swagger files once
npm run swagger:convert

# Watch for changes and convert automatically
npm run swagger:watch
```

See the example `package.json` in `scripts/examples/package.json.example`.

## Troubleshooting

### Hooks not running

If hooks are not running, check:

1. Ensure the hooks are executable:
   ```bash
   chmod +x .git/hooks/pre-commit .git/hooks/post-checkout
   ```

2. Ensure `swagger-to-http-file` is installed and in your PATH:
   ```bash
   which swagger-to-http-file
   ```

3. Check if hooks are being skipped due to the environment variable:
   ```bash
   unset SWAGGER_TO_HTTP_SKIP_HOOKS
   ```

### HTTP files not updated

If HTTP files are not being updated, check:

1. Ensure your Swagger files are valid:
   ```bash
   swagger-to-http-file -i your-file.json -v
   ```

2. Try running the conversion manually:
   ```bash
   swagger-to-http-file -i your-file.json -o output-dir -w -v
   ```

3. Check for file permission issues:
   ```bash
   ls -la output-dir
   ```

## Best Practices

1. **Include HTTP files in version control** to allow for easy API testing by team members.

2. **Use a consistent directory structure** for Swagger files and HTTP files.

3. **Use Git hooks in all development environments** to ensure consistency.

4. **Include validation in your CI/CD pipeline** to catch issues early.

5. **Document API changes** in both Swagger files and commit messages.