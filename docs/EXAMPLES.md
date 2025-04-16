# Examples and Use Cases

This document provides examples of how to use `swagger-to-http-file` for different scenarios and use cases.

## Table of Contents

- [Basic Conversion](#basic-conversion)
- [Custom Base URL](#custom-base-url)
- [Custom Output Directory](#custom-output-directory)
- [Group by Tags](#group-by-tags)
- [Overwriting Existing Files](#overwriting-existing-files)
- [Verbose Mode](#verbose-mode)
- [Authentication Examples](#authentication-examples)
- [Path Parameters](#path-parameters)
- [Query Parameters](#query-parameters)
- [Request Bodies](#request-bodies)
- [File Upload](#file-upload)
- [Complex Swagger Files](#complex-swagger-files)
- [Integration with IDEs](#integration-with-ides)
- [Using Git Hooks](#using-git-hooks)

## Basic Conversion

Convert a Swagger file to HTTP files with default settings:

```bash
swagger-to-http-file -i swagger.json
```

This will:
- Parse the Swagger file
- Create .http files in the current directory
- Group endpoints by tags into separate files
- Use the base URL from the Swagger file

## Custom Base URL

Specify a custom base URL for the API requests:

```bash
swagger-to-http-file -i swagger.json -b https://api.dev.example.com
```

This is useful for:
- Testing against different environments (dev, staging, production)
- When the base URL in the Swagger file is not up-to-date
- Local development with a different server URL

## Custom Output Directory

Save the HTTP files to a specific directory:

```bash
swagger-to-http-file -i swagger.json -o api-tests
```

This will create the output directory if it doesn't exist.

## Group by Tags

By default, endpoints are grouped by their tags into separate files. To disable this and create a single file with all endpoints:

```bash
swagger-to-http-file -i swagger.json -g=false
```

This will create a single file named `all-endpoints.http` in the output directory.

## Overwriting Existing Files

By default, the tool will not overwrite existing HTTP files. To force overwriting:

```bash
swagger-to-http-file -i swagger.json -w
```

## Verbose Mode

Enable verbose output for debugging and more information:

```bash
swagger-to-http-file -i swagger.json -v
```

## Authentication Examples

The generated HTTP files include authentication placeholders based on the security schemes defined in the Swagger file.

### Bearer Token Authentication

```
### Get User Profile
GET {{baseUrl}}/users/profile
Authorization: Bearer {{accessToken}}
```

### Basic Authentication

```
### Login
POST {{baseUrl}}/login
Authorization: Basic {{base64_credentials}}
Content-Type: application/json

{
  "username": "user",
  "password": "pass"
}
```

### API Key Authentication

```
### Get Resources
GET {{baseUrl}}/resources
X-API-Key: {{apiKey}}
```

## Path Parameters

HTTP files with path parameters are generated with placeholders:

```
### Get User by ID
GET {{baseUrl}}/users/{{userId}}
Accept: application/json
```

## Query Parameters

Query parameters are included with example values:

```
### Search Users
GET {{baseUrl}}/users?name={{name}}&age={{age}}&active={{active}}
Accept: application/json
```

## Request Bodies

Request bodies are formatted with example values from the Swagger definition:

```
### Create User
POST {{baseUrl}}/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "active": true
}
```

## File Upload

For endpoints that accept file uploads:

```
### Upload Avatar
POST {{baseUrl}}/users/{{userId}}/avatar
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary

------WebKitFormBoundary
Content-Disposition: form-data; name="avatar"; filename="avatar.png"
Content-Type: image/png

< ./avatar.png
------WebKitFormBoundary--
```

## Complex Swagger Files

For complex Swagger files with nested schemas and references:

```bash
swagger-to-http-file -i complex-api.json -v
```

The tool will resolve all references and generate appropriate HTTP requests with example values.

## Integration with IDEs

### VS Code

1. Install the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension
2. Open the generated .http files
3. Click on "Send Request" above each request to execute it

### JetBrains IDEs (IntelliJ, WebStorm, etc.)

1. Open the generated .http files
2. Use the built-in HTTP client to execute requests
3. Create environment configurations for different base URLs

## Using Git Hooks

The tool includes Git hooks for automatically updating HTTP files when Swagger files change:

```bash
# Install Git hooks
make hooks
```

See [Git Hooks Documentation](GIT_HOOKS.md) for more details.

## Example Output Files

Here's an example of the structure of generated HTTP files when using the default settings:

```
http-files/
├── pet.http        # Endpoints tagged with "pet"
├── store.http      # Endpoints tagged with "store"
└── user.http       # Endpoints tagged with "user"
```

Each file contains the relevant endpoints for that tag, with appropriate placeholders and example values.
