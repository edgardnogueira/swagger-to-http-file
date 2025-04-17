# Exemplo Básico

Este exemplo demonstra o uso mais simples da ferramenta swagger-to-http-file para converter um documento Swagger em arquivos HTTP.

## Arquivo Swagger

O arquivo [petstore.json](./petstore.json) é um documento Swagger 2.0 que descreve uma API simples de petshop com operações básicas de CRUD para entidades "Pet".

## Comandos

### Conversão Básica

```bash
swagger-to-http-file -i petstore.json
```

Este comando irá gerar arquivos .http na pasta atual, organizados por tags (por padrão).

### Conversão com Diretório de Saída Específico

```bash
swagger-to-http-file -i petstore.json -o output
```

### Conversão com Arquivo Único (Sem Organização por Tags)

```bash
swagger-to-http-file -i petstore.json -g=false
```

## Resultado

### Estrutura Gerada

```
.
├── pets.http  # Agrupado por tag "pets"
```

### Conteúdo de pets.http

```http
# Global variables
@baseUrl = http://petstore.swagger.io/api

### List all pets
GET {{baseUrl}}/pets
Accept: application/json

### Create a pet
POST {{baseUrl}}/pets
Content-Type: application/json

{
  "id": 0,
  "name": "string",
  "tag": "string"
}

### Info for a specific pet
GET {{baseUrl}}/pets/{{petId}}
Accept: application/json

### Update a pet
PUT {{baseUrl}}/pets/{{petId}}
Content-Type: application/json

{
  "id": 0,
  "name": "string",
  "tag": "string"
}

### Delete a pet
DELETE {{baseUrl}}/pets/{{petId}}
Accept: application/json
```

## Características Destacadas

- Conversão automática de endpoints para formato HTTP
- Inclusão de variáveis globais como `@baseUrl`
- Agrupamento por tags
- Adição de headers apropriados (Content-Type, Accept)
- Inclusão de exemplos de corpos para requisições POST/PUT
