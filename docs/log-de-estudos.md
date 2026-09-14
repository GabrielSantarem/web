# Log de Estudos

## DATA: 07/09/2026
**Objetivo:** Implementar um mock de criação de usuário em memória para validar o conceito de portas e adaptadores.

**Anotações:** 
O desenvolvimento exigiu a separação estrita entre a regra de negócio e a infraestrutura. O padrão demonstrou ser eficaz para o isolamento do domínio, porém introduz complexidade na gestão das interfaces. Para evitar a proliferação excessiva de contratos (boilerplate) para cada entidade de domínio, a solução identificada foi a adoção de Generics nas portas de saída (Repositórios).
