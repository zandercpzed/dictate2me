# Política de Segurança

## Versões Suportadas

| Versão | Suportada          |
| ------ | ------------------ |
| 1.x.x  | :white_check_mark: |
| < 1.0  | :x:                |

## Reportando uma Vulnerabilidade

A segurança do dictate2me é levada a sério. Se você descobriu uma vulnerabilidade de segurança, por favor, siga estas etapas:

### ⚠️ NÃO reporte vulnerabilidades via Issues públicas

### Como Reportar

1. **Email**: Envie um email para **security@dictate2me.dev**
2. **Assunto**: Use o prefixo `[SECURITY]` no assunto
3. **Conteúdo**: Inclua o máximo de detalhes possível:
   - Tipo de vulnerabilidade
   - Passos para reproduzir
   - Impacto potencial
   - Sugestões de correção (se houver)

### O que Esperar

- **Confirmação**: Responderemos em até 48 horas confirmando o recebimento
- **Avaliação**: Avaliaremos a vulnerabilidade em até 7 dias
- **Correção**: Trabalharemos em uma correção e coordenaremos a divulgação
- **Crédito**: Você será creditado na release notes (se desejar)

### Divulgação Responsável

Pedimos que:

- Nos dê tempo razoável para corrigir antes de divulgar publicamente
- Não explore a vulnerabilidade além do necessário para demonstrá-la
- Não acesse ou modifique dados de outros usuários

## Práticas de Segurança

### Assinatura de Commits

Todos os commits na branch main devem ser assinados. Verifique com:

```bash
git log --show-signature
```

### Verificação de Binários

Releases são assinados. Verifique com:

```bash
cosign verify-blob --key dictate2me.pub dictate2me-darwin-arm64.tar.gz
```

### Dependências

- Usamos Dependabot para atualizações automáticas
- Todas as dependências são verificadas via `go mod verify`
- Executamos `govulncheck` no CI

## Modelo de Ameaças

### Escopo

dictate2me processa:

- Áudio do microfone
- Texto transcrito
- Configurações do usuário

### Garantias

- ✅ Dados nunca são enviados para servidores externos
- ✅ Modelos de IA rodam 100% localmente
- ✅ Configurações são armazenadas em arquivos locais

### Limitações

- ⚠️ Não protegemos contra acesso físico ao dispositivo
- ⚠️ Logs podem conter trechos de texto transcrito
- ⚠️ Modelos de IA podem ter vieses

## Auditoria

Planejamos realizar auditorias de segurança regulares. Os relatórios serão publicados em `/docs/security/`.

## Vulnerabilidades Conhecidas

Nenhuma no momento. Este arquivo será atualizado conforme necessário.

## Histórico de Segurança

| Data | Versão | Severidade | Descrição | Status |
| ---- | ------ | ---------- | --------- | ------ |
| -    | -      | -          | -         | -      |

---

Para questões gerais de segurança que não sejam vulnerabilidades, use [GitHub Discussions](https://github.com/zandercpzed/dictate2me/discussions).
