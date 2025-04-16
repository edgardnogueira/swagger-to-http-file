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

For Node.js projects, you can add the following scripts to your `package.json`:

```json
{
  "scripts": {
    "swagger:convert": "swagger-to-http-file -i ./api/swagger.json -o ./http -w",
    "swagger:watch": "nodemon --watch ./api --ext json,yaml,yml --exec 'npm run swagger:convert'"
  }
}
```

This allows you to:
- Run `npm run swagger:convert` to manually convert files
- Run `npm run swagger:watch` to watch for changes and automatically convert files

A complete example `package.json` is available in `scripts/examples/package.json.example`.

## Troubleshooting

### Hook Not Running

- Make sure the hook files are executable: `chmod +x .git/hooks/pre-commit .git/hooks/post-checkout`
- Check if `SWAGGER_TO_HTTP_SKIP_HOOKS` is set in your environment
- Ensure `swagger-to-http-file` is installed and in your PATH

### HTTP Files Not Being Updated

- Run `swagger-to-http-file -i file.json -o output_dir -v -w` manually to see the output
- Check if the Swagger files are valid by running `swagger-to-http-file -i file.json --validate`
- Try running the detect script manually: `./scripts/detect-swagger-changes.sh --verbose`

### Conflicts with Other Hooks

If you're using other Git hook systems (like pre-commit, commitlint, etc.), you may need to integrate the Swagger-to-HTTP logic into those hooks. The Husky setup script provides a good example of how to do this.

## Advanced Usage

### Custom File Detection

By default, the hooks look for `.json`, `.yaml`, and `.yml` files. If you need to customize this:

1. Edit the hook scripts in `.git/hooks/` or the Husky scripts
2. Set the `SWAGGER_TO_HTTP_FILES` environment variable with specific files

### Multiple Output Directories

If you need to generate HTTP files in multiple directories based on the Swagger files:

1. Run the manual script separately for each directory:
   ```bash
   ./scripts/detect-swagger-changes.sh -d api/v1 -o http/v1
   ./scripts/detect-swagger-changes.sh -d api/v2 -o http/v2
   ```

2. Or customize the hook scripts to handle multiple directories