# Diário de Estudos e Decisões de Arquitetura

Registro cronológico das decisões técnicas, experimentos e evolução da arquitetura do projeto.

---

## 07/09/2026 - Setup Inicial e Conceituação de Portas & Adaptadores
**Objetivo:** Implementar um mock de criação de usuário em memória para validar o conceito de portas e adaptadores.

**Anotações:** 
O desenvolvimento exigiu a separação estrita entre a regra de negócio e a infraestrutura. O padrão demonstrou ser eficaz para o isolamento do domínio, porém introduz complexidade na gestão das interfaces. Para evitar a proliferação excessiva de contratos (boilerplate) para cada entidade de domínio, a solução identificada foi o planejamento de Generics nas portas de saída (Repositórios).

---

## 13/09/2026 - Camada HTTP, Validação e Padronização de Erros
**Objetivo:** Estruturar os adaptadores de entrada (handlers HTTP) e tratamento robusto de requisições.

**Anotações:**
* Adoção da biblioteca `validator/v10` para validação de DTOs (`UserCreateRequest`).
* Criação de `ErrorResponse` para garantir consistência nas mensagens de erro da API.
* Implementação do middleware de logging e graceful shutdown no servidor HTTP (`main.go`).
* Identificado ponto de melhoria no desacoplamento dos códigos de status HTTP em relação aos erros de domínio.

---

## 17/09/2026 - Abstração Genérica de Persistência com Generics
**Objetivo:** Reduzir boilerplate nas portas de saída aplicando Go Generics.

**Anotações:**
* Definição da interface restritiva `core.Entity` exigindo o método `GetID() string`.
* Criação da interface genérica `core.Repository[T Entity]` definindo as operações base: `Create`, `Update`, `FindByID`, `FindAll` e `Delete`.
* A entidade `User` passou a implementar `core.Entity`. Dessa forma, o compilador Go faz a monomorfização em tempo de compilação sem custos de alocação de `any` ou type assertions em tempo de execução.

---

## 19/09/2026 - Adaptador de Memória Concorrente e Testes de Corrida
**Objetivo:** Implementar o adaptador de persistência genérico em memória com proteção para acessos concorrentes.

**Anotações:**
* Criação do `MemoryRepository[T]` no pacote `adapters`.
* Como Go gerencia requisições HTTP em goroutines separadas, estruturas em memória compartilhadas precisam de sincronização.
* Aplicação de `sync.RWMutex` (com `RLock`/`RUnlock` para leituras e `Lock`/`Unlock` para escritas).
* Proteção estendida também ao `InMemoryUserRepository`.
* Execução e validação de testes de estresse concorrente validados com a flag `-race` do Go.

---

## 22/09/2026 - Testes de Domínio Isolados (UserService)
**Objetivo:** Garantir a integridade das regras de negócio do núcleo de domínio sem dependência de adaptadores externos.

**Anotações:**
* Criação de um mock em memória do `UserRepository` exclusivo para a suíte de testes de `UserService`.
* Cobertura de cenários de sucesso e caminhos de falha: validações de campos obrigatórios (`ErrInvalidInput`), busca de usuário inexistente (`ErrUserNotFound`) e ativação de conta já ativa (`ErrUserAlreadyActive`).
* Os testes comprovam a eficácia da arquitetura: o núcleo foi 100% testado sem levantar servidor HTTP ou banco de dados real.

---

## 24/09/2026 - Conclusão dos Endpoints da API REST
**Objetivo:** Expor as operações de consulta e ativação de usuário através do protocolo HTTP.

**Anotações:**
* Implementação de `HandleGetUser` (`GET /users/{id}`) e `HandleActivateUser` (`PATCH /users/{id}/activate`).
* Utilização dos novos recursos de roteamento nativo do Go 1.22 (`http.ServeMux` com casamento de método e extração direta via `r.PathValue`).
* Mapeamento semântico dos erros de domínio para códigos HTTP correspondentes:
  * `ErrUserNotFound` -> `404 Not Found`
  * `ErrEmailAlreadyUsed` / `ErrUserAlreadyActive` -> `409 Conflict`
  * `ErrValidationFailure` / `ErrInvalidInput` -> `400 Bad Request`

---

## 26/09/2026 - Conclusões para a Apresentação
**Pontos fortes observados no padrão:**
1. **Testabilidade:** O isolamento das camadas torna os testes unitários triviais e rápidos.
2. **Substituibilidade:** Trocar o repositório em memória por um banco relacional (ex: PostgreSQL com pgx) não altera uma única linha de código do `internal/core`.
3. **Clareza de responsabilidades:** A camada de transporte HTTP apenas serializa/deserializa e encaminha chamadas para o serviço de domínio.
