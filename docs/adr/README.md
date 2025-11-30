# Architecture Decision Records (ADRs)

Este diretório contém os Architecture Decision Records (ADRs) do projeto dictate2me.

## O que é um ADR?

Um ADR documenta uma decisão arquitetural significativa tomada no projeto, incluindo:

- O contexto que motivou a decisão
- A decisão em si
- Alternativas consideradas
- Consequências esperadas

## Por que usar ADRs?

- **Histórico**: Preserva o raciocínio por trás de decisões importantes
- **Conhecimento**: Novos contribuidores entendem o "porquê" das escolhas
- **Revisão**: Facilita revisitar decisões quando o contexto muda
- **Discussão**: Estrutura debates técnicos

## Como criar um novo ADR?

1. Copie o [template.md](template.md)
2. Renomeie para `NNNN-titulo-curto.md` (ex: `0002-escolha-de-banco.md`)
3. Preencha todas as seções
4. Marque status como "Proposto"
5. Abra uma Pull Request para discussão
6. Após aprovação, mude status para "Aceito"

## Lista de ADRs

| #                            | Título                             | Status    | Data       |
| ---------------------------- | ---------------------------------- | --------- | ---------- |
| [0001](0001-linguagem-go.md) | Uso de Go como Linguagem Principal | ✅ Aceito | 2025-01-30 |

## Status Possíveis

- **Proposto**: ADR em discussão
- **Aceito**: Decisão aprovada e implementada
- **Depreciado**: Ainda válido mas planejado para ser substituído
- **Substituído por ADR-XXXX**: Decisão foi revisada
- **Rejeitado**: Proposta não aprovada

## Boas Práticas

1. **Seja claro e objetivo**: ADRs devem ser fáceis de ler
2. **Documente alternativas**: Mostre que considerou outras opções
3. **Seja honesto sobre consequências**: Toda decisão tem trade-offs
4. **Inclua referências**: Links para docs, benchmarks, discussões
5. **Mantenha atualizado**: Se uma decisão muda, crie novo ADR

## Quando criar um ADR?

Crie um ADR quando:

- Escolher tecnologias principais (linguagem, framework, banco de dados)
- Definir padrões arquiteturais (monolito vs microservices, REST vs GraphQL)
- Fazer mudanças que afetam múltiplos módulos
- Escolher entre múltiplas abordagens viáveis

Não precisa ADR para:

- Decisões triviais ou óbvias
- Detalhes de implementação locais
- Mudanças cosméticas

## Referências

- [ADR GitHub Organization](https://adr.github.io/)
- [Documenting Architecture Decisions - Michael Nygard](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
