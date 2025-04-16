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

## Integration with npm Projects

For Node.js projects, you can use the provided example `package.json` configuration:

1. Copy the relevant parts from `scripts/examples/package.json.example` to your `package.json`
2. Install the necessary dependencies:
   ```bash
   npm install --save-dev husky nodemon
   ```

3. This adds the following npm scripts:
   - `npm run swagger:convert`: Run a one-time conversion
   - `npm run swagger:watch`: Watch for changes and convert automatically

## Troubleshooting

### Hooks Not Running

1. Make sure the hooks are executable:
   ```bash
   chmod +x .git/hooks/pre-commit .git/hooks/post-checkout
   ```

2. Check if hooks are being skipped:
   ```bash
   # Make sure this environment variable is not set
   unset SWAGGER_TO_HTTP_SKIP_HOOKS
   ```

3. Verify the tool is installed and in your PATH:
   ```bash
   which swagger-to-http-file
   ```

### HTTP Files Not Updated

1. Run the detection script with verbose output:
   ```bash
   ./scripts/detect-swagger-changes.sh -v
   ```

2. Check if your Swagger files are detected:
   ```bash
   find . -name "*.json" -o -name "*.yaml" -o -name "*.yml"
   ```

3. Try running the conversion manually:
   ```bash
   swagger-to-http-file -i path/to/swagger.json -o output/dir -w -v
   ```

## Best Practices

1. **Commit Swagger and HTTP files together**: This maintains a clear history of API changes.

2. **Include HTTP files in version control**: This allows team members to use them without installing the tool.

3. **Configure output directories consistently**: Use the same output directory across your team.

4. **Set up CI checks**: Ensure HTTP files are always in sync with Swagger files.

5. **Document your setup**: Include a note in your project README about the Git hooks setup.