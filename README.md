# Projeto Golang API com Cache Redis

Este projeto é uma API construída em [Golang](https://golang.org/) que utiliza o framework [Gin](https://gin-gonic.com/) para roteamento HTTP e o banco de dados em memória [Redis](https://redis.io/) para cache de requisições. A API faz requisições para a API da Omie, armazena as respostas no Redis por um tempo determinado e retorna os dados de forma eficiente, utilizando cache para evitar chamadas repetidas à API externa.

## Funcionalidades

- Faz requisições para a API da Omie.
- Armazena as respostas no cache Redis por 60 segundos.
- Se os dados já estiverem no cache, retorna diretamente do Redis.
- Implementação de um sistema de mutex para garantir que múltiplas requisições sejam processadas de forma síncrona.
- Suporte ao uso de variáveis de ambiente para chaves de API e configurações do Redis.

## Tecnologias Utilizadas

- **Golang**: Linguagem de programação principal.
- **Gin**: Framework web para roteamento de requisições.
- **Redis**: Banco de dados em memória para caching.
- **go-redis**: Cliente Redis para Go.
- **godotenv**: Carregar variáveis de ambiente a partir de arquivos `.env`.

## Requisitos

- [Golang](https://golang.org/dl/) 1.17 ou superior.
- Redis instalado (pode ser via [WSL para Windows](https://docs.microsoft.com/en-us/windows/wsl/install) ou nativamente no Linux/Mac).
- **go modules** habilitado (caso o projeto já use `go.mod` e `go.sum`).

## Configuração

### Instalação

1. Clone o repositório:
```bash
   git clone https://github.com/seu-usuario/seu-repositorio.git
```
2. Navegue até o diretório do projeto:
```bash
   cd seu-repositorio
```
3. Instale as dependências do projeto:
```bash
   go mod tidy
```
4. Crie um arquivo .env na raiz do projeto com as variáveis de ambiente necessárias para o Redis e a API Omie:
```makefile
   OMIE_APP_KEY=your_omie_app_key
   OMIE_APP_SECRET=your_omie_app_secret
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
```
## Executar o Projeto
Certifique-se de que o Redis está rodando. Se estiver utilizando o WSL, você pode iniciar o Redis com o comando:
1. Certifique-se de que o Redis está rodando. Se estiver utilizando o WSL, você pode iniciar o Redis com o comando:
```bash
   sudo service redis-server start
```
2. Execute a aplicação Golang:
```bash
   go run main.go
```
3. A API estará rodando em: http://localhost:8080.
## Testando a API
Faça uma requisição POST para o endpoint /omie_request com o seguinte corpo:
```json
{
  "call": "nome_da_funcao_omie",
  "param": [
    {
      "chave1": "valor1",
      "chave2": "valor2"
    }
  ],
  "url": "caminho/da/api/omie"
}
```
Por exemplo, usando `curl`:
```bash
curl -X POST http://localhost:8080/omie_request \
-H "Content-Type: application/json" \
-d '{
    "call": "ListarClientes",
    "param": [{"pagina": 1}],
    "url": "geral/clientes/"
}'
```
Se os dados já estiverem no cache, eles serão retornados a partir dele. Caso contrário, a API Omie será consultada, e os dados serão armazenados no cache Redis.

## Licença
Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

### Explicação:
- **Configuração do ambiente**: Instruções sobre como configurar variáveis de ambiente e instalar dependências.
- **Uso da API**: Exemplos de requisição POST para testar o funcionamento da API com Redis cache.
- **Script Redis**: Referência ao script PowerShell de instalação do Redis, facilitando a instalação no Windows.

