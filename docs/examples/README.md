# Exemplos de Uso do Swagger-to-HTTP-File

Esta pasta contém exemplos práticos de como usar a ferramenta swagger-to-http-file em diferentes cenários.

## Estrutura

- **[basic/](./basic/)** - Exemplo básico mostrando conversão simples de Swagger para HTTP
- **[authentication/](./authentication/)** - Exemplos com diferentes métodos de autenticação
- **[parameters/](./parameters/)** - Exemplos de uso com diferentes tipos de parâmetros
- **[organization/](./organization/)** - Exemplos de organização por tags/diretórios

## Como Usar os Exemplos

Cada pasta de exemplo contém:
1. Um arquivo Swagger JSON/YAML
2. Os arquivos HTTP correspondentes gerados
3. Um README explicando o cenário

Para testar qualquer exemplo:

```bash
# Navegue até a pasta do exemplo
cd docs/examples/basic

# Execute a ferramenta com o arquivo Swagger de exemplo
swagger-to-http-file -i petstore.json -o output

# Compare o resultado com os arquivos .http fornecidos
```
