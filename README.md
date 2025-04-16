# Swagger to HTTP File Converter

A command-line tool that converts Swagger/OpenAPI JSON documents into `.http` files for easy API testing.

## Features

- Parse Swagger/OpenAPI JSON files
- Generate `.http` files with proper formatting
- Organize requests by tags/directories
- Support for path, query, and body parameters
- Support for authentication mechanisms
- Group requests by tags into separate files
- Git hooks for automatic HTTP file updates when Swagger files change

## Installation

See [Installation Guide](docs/INSTALL.md) for detailed instructions.

Quick install with Go:

```bash
go install github.com/edgardnogueira/swagger-to-http-file/cmd/swagger-to-http-file@latest
```

## Usage

### Basic Usage

Convert a Swagger file to HTTP files:

```bash
swagger-to-http-file -i swagger.json
```

This will create `.http` files in the current directory, organized by tags.

### Command-line Options

```
Usage:
  swagger-to-http-file [flags]
  swagger-to-http-file [command]

Available Commands:
  help        Help about any command
  version     Print the version information

Flags:
  -b, --baseUrl string     Base URL for API requests (overrides the one in Swagger)
  -g, --group-by-tag       Group requests by tags into separate files (default true)
  -h, --help               help for swagger-to-http-file
  -i, --input string       Swagger/OpenAPI JSON file to convert (required)
  -o, --output string      Directory to save .http files (default ".")
  -w, --overwrite          Overwrite existing files
  -v, --verbose            Enable verbose output
```

### Examples

Convert a Swagger file and output to a specific directory:

```bash
swagger-to-http-file -i swagger.json -o http-requests
```

Convert a Swagger file with a custom base URL:

```bash
swagger-to-http-file -i swagger.json -b https://api.example.com
```

Create a single file with all requests:

```bash
swagger-to-http-file -i swagger.json -g=false
```

Override existing files:

```bash
swagger-to-http-file -i swagger.json -w
```

## HTTP File Format

The generated `.http` files follow the format recognized by tools like VS Code's REST Client extension or JetBrains IDEs. Example:

```
# Global variables
@baseUrl = https://api.example.com
@authToken = your_auth_token

### Get Pets
GET {{baseUrl}}/pets
Accept: application/json

### Create Pet
POST {{baseUrl}}/pets
Content-Type: application/json

{
  "name": "string",
  "age": 0
}

### Get Pet by ID
GET {{baseUrl}}/pets/{{petId}}
Accept: application/json
```

## Git Hooks Integration

The tool provides Git hooks for automatically updating HTTP files when Swagger/OpenAPI files change:

- **Pre-commit hook**: Automatically updates HTTP files when Swagger files are committed
- **Post-checkout hook**: Updates HTTP files when switching branches

To install the Git hooks:

```bash
# Make the script executable
chmod +x scripts/install-hooks.sh

# Run the installation script
./scripts/install-hooks.sh
```

For more details, see the [Git Hooks Documentation](docs/GIT_HOOKS.md).

## Development

This project follows Clean Architecture principles and is developed in Go.

### Project Structure

```
swagger-to-http-file/
├── cmd/                      # Command-line entry points
│   └── swagger-to-http-file/ # Main application
├── internal/                 # Private application code
│   ├── domain/               # Domain models
│   ├── application/          # Application layer
│   ├── adapters/             # Adapter layer
│   └── infrastructure/       # Infrastructure layer
├── docs/                     # Documentation
└── test/                     # Test files and samples
```

### Building from Source

The project uses a Makefile to simplify common development tasks:

```bash
# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Check for linting issues
make lint

# Run with an example input
make run

# Install Git hooks
make hooks
```

### Docker

You can build and run the project with Docker:

```bash
# Build the Docker image
docker build -t swagger-to-http-file .

# Run the container
docker run --rm -v $(pwd):/data swagger-to-http-file -i /data/swagger.json -o /data
```

### CI/CD

The project uses GitHub Actions for:
- Running tests and linting on every pull request
- Building and releasing binaries for multiple platforms
- Publishing Docker images

See the [Contributing Guide](CONTRIBUTING.md) for information on how to contribute to the project.

## License

[MIT](LICENSE)
