# Generics em Go

A utilização de Generics permite a criação de estruturas que operam de forma independente de um tipo específico de dado, eliminando a redundância estrutural sem comprometer a validação estática do compilador.

### Instanciação em Tempo de Compilação

O processo de compilação gera implementações estáticas exclusivas para cada tipo utilizado no sistema, evitando penalidades de desempenho em tempo de execução.

```mermaid
flowchart TD
    A[Contrato Genérico Base <br> Repository T] --> B{Processo de Compilação <br> Monomorfização}
    B -->|Instanciação com User| C[Repository User <br> Otimizado para User]
    B -->|Instanciação com Product| D[Repository Product <br> Otimizado para Product]
```
A Relação com Interfaces (Restrições)

Em Go, as interfaces operam como restrições (bounds) para os tipos genéricos. Elas estabelecem um filtro estrito: o parâmetro genérico só aceitará tipos que assinarem o contrato definido pela interface.
```mermaid
flowchart TD
    A[Interface Restritiva <br> Regra: T deve possuir o método GetID] --> B{Validação do Compilador}
    
    C[Entidade: User <br> Possui o método GetID] --> B
    D[Entidade: Invoice <br> Não possui o método GetID] --> B
    
    B -->|Avaliação de User| E[APROVADO <br> Tipo aceito no Genérico]
    B -->|Avaliação de Invoice| F[REJEITADO <br> Erro de compilação]
    
    style E stroke:#198754,stroke-width:2px
    style F stroke:#dc3545,stroke-width:2px
```