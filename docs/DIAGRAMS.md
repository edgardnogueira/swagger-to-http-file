# Diagrams and Visual Guides

This document provides visual representations of the tool's architecture, workflow, and functionality to help users better understand how `swagger-to-http-file` works.

## Conversion Workflow

```mermaid
graph TD
    A[Swagger/OpenAPI File] -->|Parse| B[Internal Representation]
    B -->|Process| C[Generate HTTP Files]
    C -->|Group by Tags| D[Tag-based HTTP Files]
    C -->|Single File| E[All Endpoints HTTP File]
    F[Command Line Options] -->|Configure| B
    G[Base URL] -->|Override| B
```

## Architecture

```mermaid
graph TD
    subgraph "Command Line Interface"
        CLI[CLI Entry Point]
        Parser[Flag Parser]
        Validator[Input Validator]
    end
    
    subgraph "Core Processing"
        SwaggerParser[Swagger Parser]
        OpenAPIParser[OpenAPI Parser]
        HTTPGenerator[HTTP File Generator]
    end
    
    subgraph "Output Management"
        FileWriter[File Writer]
        DirectoryManager[Directory Manager]
    end
    
    CLI --> Parser
    Parser --> Validator
    Validator --> SwaggerParser
    Validator --> OpenAPIParser
    SwaggerParser --> HTTPGenerator
    OpenAPIParser --> HTTPGenerator
    HTTPGenerator --> FileWriter
    FileWriter --> DirectoryManager
```

## Directory Structure

```mermaid
graph TD
    subgraph "Repository Structure"
        Root[swagger-to-http-file]
        Cmd[cmd/]
        Internal[internal/]
        Docs[docs/]
        Test[test/]
        Scripts[scripts/]
    end
    
    Root --> Cmd
    Root --> Internal
    Root --> Docs
    Root --> Test
    Root --> Scripts
    
    subgraph "cmd/ directory"
        CmdMain[swagger-to-http-file/]
        CmdMain --> Main[main.go]
    end
    
    subgraph "internal/ directory"
        Domain[domain/]
        Application[application/]
        Adapters[adapters/]
        Infrastructure[infrastructure/]
    end
    
    Internal --> Domain
    Internal --> Application
    Internal --> Adapters
    Internal --> Infrastructure
    
    subgraph "Clean Architecture"
        DomainModels[Domain Models]
        UseCases[Use Cases]
        InterfaceAdapters[Interface Adapters]
        Frameworks[Frameworks & Drivers]
    end
    
    Domain --> DomainModels
    Application --> UseCases
    Adapters --> InterfaceAdapters
    Infrastructure --> Frameworks
```

## Git Hooks Integration

```mermaid
sequenceDiagram
    participant D as Developer
    participant PCS as Pre-Commit Script
    participant S2H as swagger-to-http-file
    participant G as Git
    
    D->>G: git commit (with swagger file changes)
    G->>PCS: Trigger pre-commit hook
    PCS->>S2H: Detect Swagger file changes
    S2H->>S2H: Generate HTTP files
    S2H->>PCS: Return generated files
    PCS->>G: Add generated files to commit
    G->>D: Complete commit
```

## HTTP File Generation Process

```mermaid
graph TD
    A[Parse Swagger/OpenAPI File] --> B{Valid File?}
    B -->|Yes| C[Process API Endpoints]
    B -->|No| Error[Return Error]
    C --> D[For Each Endpoint]
    D --> E[Extract Method, Path, Parameters]
    D --> F[Extract Request Body]
    D --> G[Extract Response Types]
    D --> H[Format HTTP Request]
    H --> I{Group by Tags?}
    I -->|Yes| J[Create File per Tag]
    I -->|No| K[Create Single File]
    J --> L[Write Files to Output Directory]
    K --> L
```

## Command Line Usage

```mermaid
graph LR
    A[swagger-to-http-file] -->|required| B["-i swagger.json"]
    A --> C["-o output_dir"]
    A --> D["-b custom_base_url"]
    A --> E["-g=false"]
    A --> F["-w"]
    A --> G["-v"]
```

## HTTP File Structure

```
# Global variables
@baseUrl = https://api.example.com
@apiKey = your_api_key

### Get Pet by ID
GET {{baseUrl}}/pets/{{petId}}
Accept: application/json

### Create Pet
POST {{baseUrl}}/pets
Content-Type: application/json

{
  "name": "Fluffy",
  "species": "cat",
  "age": 3
}
```

## IDE Integration

### VS Code REST Client

After generating HTTP files, you can use them in VS Code with the REST Client extension:

![VS Code REST Client](https://raw.githubusercontent.com/Huachao/vscode-restclient/master/images/usage.gif)

### JetBrains HTTP Client

Similarly, in JetBrains IDEs, the built-in HTTP client can execute the generated requests:

![JetBrains HTTP Client](https://www.jetbrains.com/help/img/idea/2023.1/http_response_example.png)

## Complete Workflow

```mermaid
graph TD
    A[Swagger/OpenAPI Document] -->|Develop API| B[API Development]
    B -->|Generate Documentation| A
    A -->|Convert| C[swagger-to-http-file]
    C -->|Generate| D[HTTP Files]
    D -->|Test in IDE| E[API Testing]
    E -->|Feedback| B
    
    F[Git Hooks] -->|Automate| C
    G[Manual CLI] -->|Run| C
    
    H[CI/CD Pipeline] -->|Verify| D
    D -->|Version Control| I[Git Repository]
    
    J[Developers] -->|Use| D
    J -->|Update| A
```

These diagrams provide a visual understanding of the tool's architecture, workflow, and usage patterns. They can be rendered by opening this markdown file in a viewer that supports Mermaid diagrams, such as GitHub, or using tools like Mermaid Live Editor.
