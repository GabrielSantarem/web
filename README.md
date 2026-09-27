# Estudo de Arquitetura Hexagonal em Go

Projeto simples que fiz para estudar arquitetura hexagonal (portas e adaptadores) e generics no Go.

## Estrutura das pastas

- `internal/core`: regras de negocio e interfaces (portas).
- `internal/adapters`: repositorio em memoria (adaptador de saida).
- `handler`: handlers http e validacao dos dados (adaptador de entrada).
- `docs/`: anotacoes de estudo.

## Como rodar

Rodar a aplicacao:
```bash
go run main.go
```

Rodar os testes:
```bash
go test ./...
```

## Rotas da API

- `GET /`: rota simples de hello.
- `POST /users/create`: cria usuario (envia json com `name` e `email`).
- `GET /users/{id}`: busca um usuario pelo id.
- `PATCH /users/{id}/activate`: ativa um usuario.

## Documentos

- [Diario de Estudos](./docs/log-de-estudos.md)
- [Arquitetura Hexagonal](./docs/arquitetura-hexagonal.md)
- [Generics no Go](./docs/go-generics.md)
