# Installation Guide

This document describes how to install and use the `swagger-to-http-file` tool.

## Installation Options

### Using Go Install

If you have Go installed (version 1.16+), you can install directly from the source:

```bash
go install github.com/edgardnogueira/swagger-to-http-file/cmd/swagger-to-http-file@latest
```

This will install the latest version of the tool in your `$GOPATH/bin` directory.

### Downloading Pre-built Binaries

Pre-built binaries for various platforms will be available in the GitHub releases section.

1. Go to [Releases](https://github.com/edgardnogueira/swagger-to-http-file/releases)
2. Download the appropriate binary for your operating system
3. Extract the binary (if needed)
4. Make it executable (if on Unix-like systems):
   ```bash
   chmod +x swagger-to-http-file
   ```
5. Move it to a directory in your PATH or use it from the current location

### Building from Source

To build from source:

1. Clone the repository:
   ```bash
   git clone https://github.com/edgardnogueira/swagger-to-http-file.git
   ```
2. Navigate to the project directory:
   ```bash
   cd swagger-to-http-file
   ```
3. Build the project:
   ```bash
   go build -o swagger-to-http-file ./cmd/swagger-to-http-file
   ```
4. Run the binary:
   ```bash
   ./swagger-to-http-file --help
   ```

## Verifying Installation

After installation, verify the tool is working correctly:

```bash
swagger-to-http-file version
```

This should display the version information of the tool.

## Usage

See the [README.md](../README.md) for detailed usage instructions.

## Troubleshooting

### Common Issues

- **Command not found**: Ensure the binary is in your PATH
- **Permission denied**: Make sure the binary is executable (`chmod +x swagger-to-http-file`)
- **Failed to parse Swagger file**: Check that your Swagger file is valid JSON/YAML

If you encounter any other issues, please [file an issue](https://github.com/edgardnogueira/swagger-to-http-file/issues) on GitHub.
