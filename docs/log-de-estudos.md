# Diario de Estudos

Anotacoes curtas sobre o que foi feito e decidido em cada etapa do projeto.

---

## 07/09/2026 - Inicio do projeto e Portas e Adaptadores

- Criacao da estrutura base separando `core` e `adapters`.
- Implementacao inicial de `User` e `InMemoryUserRepository`.
- Percebi que criar uma interface de repositorio para cada struct gera muito codigo repetido. Decidi estudar Generics para resolver isso.

---

## 13/09/2026 - Handlers HTTP e validacao

- Adicionado `http.ServeMux` com rotas basicas no Go 1.22.
- Adicionada validacao de dados com `validator/v10` no cadastro de usuario.
- Padronizacao das respostas de erro com a struct `ErrorResponse`.
- Adicionado middleware de log e graceful shutdown no `main.go`.

---

## 17/09/2026 - Criando contrato generico com Generics

- Criada a interface `core.Entity` com o metodo `GetID() string`.
- Criada a interface generica `core.Repository[T Entity]` com metodos padrao de CRUD.
- Feita a struct `User` implementar `GetID()`.

---

## 19/09/2026 - Repositorio em memoria concorrente

- Criado `MemoryRepository[T]` no pacote `adapters`.
- Como Go trata requisicoes http em goroutines diferentes, adicionei `sync.RWMutex` para evitar race condition no acesso aos dados.
- Adicionei testes com `go test -race` para garantir que nao da erro de concorrencia.

---

## 22/09/2026 - Testes de unidade no Core

- Escrevi os testes de negocio para `UserService` usando um mock de repositorio.
- Cobri os casos de criacao com dados invalidos, busca de usuario inexistente e ativacao de usuario.

---

## 24/09/2026 - Rotas de busca e ativacao

- Implementado `HandleGetUser` (`GET /users/{id}`) e `HandleActivateUser` (`PATCH /users/{id}/activate`).
- Uso de `r.PathValue("id")` do Go nativo.
- Ajuste dos status code: 404 quando nao acha, 409 quando usuario ja esta ativo ou email duplicado, e 400 para dados invalidos.

---

## 26/09/2026 - Conclusoes do estudo

- **O que valeu a pena**: isolar o dominio facilita muito escrever testes unitarios rapidos e trocar a implementacao do banco sem mexer na regra de negocio.
- **Dificuldade**: no inicio parece excesso de arquivos e interfaces para um crud simples, mas fica mais organizado conforme o projeto cresce.

---

## 27/09/2026 - Novo adaptador de saida: SQLite

- Criado `SQLiteUserRepository` implementando a mesma interface `core.UserRepository`.
- Usado o driver pure Go `modernc.org/sqlite` sem necessidade de CGO.
- Mapeamento dos erros do SQLite para os erros de negocio do `core` (ex: erro de constraint UNIQUE virando `core.ErrEmailAlreadyUsed`).
- Validado na pratica o maior beneficio de portas e adaptadores: a troca do adaptador em memoria por um banco real nao exigiu nenhuma alteracao no dominio (`core`) ou nas rotas (`handler`).
