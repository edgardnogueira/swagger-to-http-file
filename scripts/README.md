# Scripts Directory

This directory contains utility scripts for working with the Swagger to HTTP File Converter.

## Available Scripts

### Git Hooks Installation

- `install-hooks.sh`: Installs native Git hooks for automatic HTTP file generation
- `setup-husky.js`: Sets up Husky for Node.js projects to manage Git hooks

### Git Hooks

The `hooks` directory contains the Git hooks that are installed by the installation scripts:

- `pre-commit`: Runs before a commit to update HTTP files from staged Swagger files
- `post-checkout`: Runs after checking out a branch to update HTTP files based on changes

### Utilities

- `detect-swagger-changes.sh`: Detects changes in Swagger files and updates HTTP files

## Examples

The `examples` directory contains example configurations:

- `package.json.example`: Example package.json with Swagger-to-HTTP integration

## How to Use

See the [Git Hooks Documentation](../docs/GIT_HOOKS.md) for detailed instructions on how to set up and use these scripts.

## Quick Start

To install native Git hooks:

```bash
chmod +x scripts/install-hooks.sh
./scripts/install-hooks.sh
```

To set up Husky for Node.js projects:

```bash
chmod +x scripts/setup-husky.js
node scripts/setup-husky.js
```

To manually detect and update HTTP files:

```bash
chmod +x scripts/detect-swagger-changes.sh
./scripts/detect-swagger-changes.sh
```
