# CLI Reference

This document provides a comprehensive reference for all command-line options available in the `swagger-to-http-file` tool.

## Command Overview

```
swagger-to-http-file [flags]
swagger-to-http-file [command]
```

## Available Commands

| Command | Description |
|---------|-------------|
| `help`  | Help about any command |
| `version` | Print the version information |

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--baseUrl`, `-b` | `-b` | string | from Swagger | Base URL for API requests (overrides the one in Swagger) |
| `--group-by-tag`, `-g` | `-g` | boolean | `true` | Group requests by tags into separate files |
| `--help`, `-h` | `-h` | - | - | Help for swagger-to-http-file |
| `--input`, `-i` | `-i` | string | - | Swagger/OpenAPI JSON file to convert (required) |
| `--output`, `-o` | `-o` | string | `.` (current directory) | Directory to save .http files |
| `--overwrite`, `-w` | `-w` | boolean | `false` | Overwrite existing files |
| `--verbose`, `-v` | `-v` | boolean | `false` | Enable verbose output |

## Detailed Flag Descriptions

### `--baseUrl`, `-b`

Specifies the base URL to use for all API requests in the generated HTTP files. If not provided, the tool will use the base URL defined in the Swagger document.

**Example:**
```bash
swagger-to-http-file -i swagger.json -b https://api.example.com
```

This is useful for testing against different environments (development, staging, production) or when the base URL in the Swagger file is not correct for your current needs.

### `--group-by-tag`, `-g`

Controls whether the tool should create separate HTTP files for each tag in the Swagger document. By default, this is set to `true`.

**Example (create separate files):**
```bash
swagger-to-http-file -i swagger.json -g=true
```

**Example (create a single file with all endpoints):**
```bash
swagger-to-http-file -i swagger.json -g=false
```

When set to `true`, you'll get files like `pet.http`, `store.http`, etc. When set to `false`, you'll get a single file named `all-endpoints.http`.

### `--help`, `-h`

Displays help information about the command or a specific flag.

**Example (general help):**
```bash
swagger-to-http-file --help
```

**Example (help for a specific command):**
```bash
swagger-to-http-file version --help
```

### `--input`, `-i`

Specifies the input Swagger/OpenAPI JSON file to convert. This flag is required.

**Example:**
```bash
swagger-to-http-file -i swagger.json
```

The tool supports:
- JSON Swagger/OpenAPI files (both 2.0 and 3.0)
- YAML Swagger/OpenAPI files (both 2.0 and 3.0)

### `--output`, `-o`

Specifies the output directory where the HTTP files will be saved. If not provided, files are saved in the current directory.

**Example:**
```bash
swagger-to-http-file -i swagger.json -o http-requests
```

If the specified directory doesn't exist, it will be created.

### `--overwrite`, `-w`

If set, the tool will overwrite any existing HTTP files in the output directory. By default, this is set to `false`, meaning the tool will not overwrite existing files.

**Example:**
```bash
swagger-to-http-file -i swagger.json -w
```

This is useful when you want to regenerate HTTP files after making changes to the Swagger document.

### `--verbose`, `-v`

Enables verbose output, which includes more detailed information about the conversion process.

**Example:**
```bash
swagger-to-http-file -i swagger.json -v
```

Verbose output includes:
- Detailed parsing information
- Number of endpoints processed
- Errors and warnings
- Output file locations

## Environment Variables

The tool also supports the following environment variables:

| Variable | Description |
|----------|-------------|
| `SWAGGER_TO_HTTP_SKIP_HOOKS=1` | Skip running Git hooks |
| `SWAGGER_TO_HTTP_FILES="file1.json file2.yaml"` | Specific Swagger/OpenAPI files to check |
| `SWAGGER_TO_HTTP_OUTPUT_DIR="./http"` | Directory for generated HTTP files |

These environment variables are especially useful when running the tool as part of automated scripts or when setting up Git hooks.

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | General error (file not found, parsing error, etc.) |
| `2` | Invalid command-line arguments |

## Examples

### Basic Usage

```bash
swagger-to-http-file -i swagger.json
```

### Advanced Usage

```bash
swagger-to-http-file -i swagger.json -o api-tests -b https://api.dev.example.com -g=false -w -v
```

This command:
1. Reads the Swagger file `swagger.json`
2. Uses `https://api.dev.example.com` as the base URL
3. Creates a single file (not grouped by tags)
4. Saves it to the `api-tests` directory
5. Overwrites any existing files
6. Shows verbose output

## See Also

- [Quick Start Guide](QUICKSTART.md)
- [Examples and Use Cases](EXAMPLES.md)
- [Git Hooks Documentation](GIT_HOOKS.md)
