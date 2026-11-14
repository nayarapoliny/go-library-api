# 📚 Go Library Management API

Sistema de Gestão de Biblioteca em **Go (Golang)** com arquitetura em camadas, **Injeção de Dependência**, filtros avançados de pesquisa e suporte híbrido para **PostgreSQL** e armazenamento em memória Thread-Safe.

## 🚀 Funcionalidades & Arquitetura

- **Pesquisa Multi-Filtro:** Busca livros combinando `title`, `author`, `genre`, `publisher` e `year`.
- **Administração e Geração Automática de ID:** Endpoint administrativo `POST /api/v1/books` que valida regras de negócio e atribui IDs sequenciais (`BOOK-0001`).
- **Injeção de Dependência (Clean Architecture):** A camada de serviço (`Service`) depende apenas da interface `Repository`, permitindo alternar entre PostgreSQL e Memória sem alterar regras de negócio.
- **Graceful Shutdown:** Encerramento controlado de conexões ativas via `os/signal` e `context.Context`.

## 🛠️ Como Executar e Testar

    go test -v -race ./...
    go run ./cmd/api

### Exemplos de Requisições (cURL)

**1. Listar e filtrar livros por gênero e ano:**

    curl "http://localhost:8080/api/v1/books?genre=Technology&year=2017"

**2. Cadastrar um novo livro (Admin):**

    curl -X POST http://localhost:8080/api/v1/books -H "Content-Type: application/json" -d '{"title":"Domain-Driven Design","author":"Eric Evans","genre":"Software Engineering","published_year":2003,"publisher":"Addison-Wesley"}'
