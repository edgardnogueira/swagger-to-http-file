# Quick Start Guide

This guide will help you get up and running with `swagger-to-http-file` in just a few minutes.

## Installation

### Using Go

```bash
go install github.com/edgardnogueira/swagger-to-http-file/cmd/swagger-to-http-file@latest
```

### Using Pre-built Binaries

1. Download the appropriate binary for your system from the [Releases page](https://github.com/edgardnogueira/swagger-to-http-file/releases)
2. Extract it (if needed)
3. Make it executable (on Unix-like systems):
   ```bash
   chmod +x swagger-to-http-file
   ```
4. Move it to a directory in your PATH or use it from the current location

## Basic Usage

1. Navigate to a directory containing your Swagger/OpenAPI JSON file
2. Run the converter:
   ```bash
   swagger-to-http-file -i swagger.json
   ```
3. Find the generated .http files in the current directory

## Example Workflow

Let's walk through a complete example using the Petstore Swagger file:

1. Save the Petstore Swagger file:

   ```bash
   # Copy the sample file
   cp $(go env GOPATH)/src/github.com/edgardnogueira/swagger-to-http-file/test/samples/petstore.json .
   ```

   Or download it:
   
   ```bash
   curl -o petstore.json https://raw.githubusercontent.com/edgardnogueira/swagger-to-http-file/main/test/samples/petstore.json
   ```

2. Convert it to HTTP files:

   ```bash
   swagger-to-http-file -i petstore.json -o http-requests
   ```

3. Explore the generated files:

   ```bash
   ls -la http-requests/
   ```

   You should see files like `pet.http`, `store.http`, and `user.http`.

4. Open the files in an IDE with HTTP client support (VS Code with REST Client extension or JetBrains IDEs)

5. Execute requests directly from the editor

## Common Command Options

```bash
# Specify a custom base URL
swagger-to-http-file -i swagger.json -b https://api.dev.example.com

# Save to a specific directory
swagger-to-http-file -i swagger.json -o api-tests

# Create a single file instead of grouping by tags
swagger-to-http-file -i swagger.json -g=false

# Force overwriting existing files
swagger-to-http-file -i swagger.json -w

# Enable verbose output
swagger-to-http-file -i swagger.json -v
```

## Setting Up Git Hooks

To automatically update HTTP files when Swagger files change:

```bash
# Make the script executable
chmod +x scripts/install-hooks.sh

# Run the installation script
./scripts/install-hooks.sh
```

## Next Steps

- Read the [Examples and Use Cases](EXAMPLES.md) for more advanced usage
- Check the [CLI Reference](CLI_REFERENCE.md) for all available options
- See the [Git Hooks Documentation](GIT_HOOKS.md) for automating the process
