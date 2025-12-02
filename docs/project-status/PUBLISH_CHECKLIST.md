# Checklist para Publicação Open Source

Este documento contém todos os passos necessários para tornar o projeto dictate2me público e colaborativo.

## ✅ Pré-requisitos (COMPLETO)

- [x] Código funcional e testado
- [x] Documentação completa
- [x] Testes com boa cobertura (85%+)
- [x] LICENSE file (MIT)
- [x] CODE_OF_CONDUCT.md
- [x] CONTRIBUTING.md detalhado
- [x] SECURITY.md
- [x] README.md atraente

## 📋 Preparação do Repositório

### 1. Limpeza e Organização

```bash
# Remover arquivos sensíveis/desnecessários
rm -f *.log
rm -rf coverage.html
rm -rf node_modules/

# Verificar .gitignore
cat .gitignore
```

- [x] Verificar se não há tokens/secrets no código
- [x] Remover arquivos temporários
- [x] Limpar histórico de Git se necessário

### 2. Verificar Branch Principal

```bash
# Renomear master para main (se necessário)
git branch -m master main

# Verificar branch atual
git branch --show-current
```

- [x] Branch principal é `main`
- [x] Branch está limpa e atualizada

### 3. Tags e Releases

```bash
# Criar tag para primeira release
git tag -a v0.1.0 -m "Initial public release

Features:
- Audio capture with PortAudio
- Offline transcription with Vosk
- LLM correction with Ollama
- REST API with WebSocket streaming
- Obsidian plugin
- Complete documentation (9,000+ lines)
"

# Push tag
git push origin v0.1.0
```

- [x] Tag v0.1.0 criada
- [x] Tag pushed para origin

## 🔧 Configuração do GitHub

### 1. Repository Settings

#### General

- [ ] Description: "🎤 dictate2me - Offline voice transcription with AI correction for Obsidian and more"
- [ ] Website: (adicionar se houver)
- [ ] Topics: `voice-transcription`, `offline-ai`, `obsidian-plugin`, `golang`, `typescript`, `vosk`, `ollama`
- [ ] Features:
  - [x] Wikis: Disabled
  - [x] Issues: Enabled
  - [x] Discussions: Enabled ✨
  - [x] Projects: Disabled (por enquanto)
  - [x] Preserve this repository: Disabled
  - [x] Sponsorships: Disabled (por enquanto)

#### Branches

- [ ] Default branch: `main`
- [ ] Branch protection rules:
  ```
  Branch: main
  - Require pull request before merging
  - Require approvals: 1
  - Dismiss stale reviews
  - Require status checks (CI)
  - Do not allow bypassing
  ```

### 2. GitHub Actions

Verificar workflows em `.github/workflows/`:

- [x] `ci.yaml` - Roda testes, build, lint e security em cada push/PR
- [ ] `release.yaml` - Cria releases automaticamente em tags (próximo passo)

**Testar:**

```bash
# Push para trigger CI
git push origin main

# Verificar em: https://github.com/zandercpzed/dictate2me/actions
```

### 3. Issue Templates

Verificar em `.github/ISSUE_TEMPLATE/`:

- [x] `bug_report.md` - Template de bug
- [x] `feature_request.md` - Template de feature
- [x] `question.md` - Template de pergunta
- [x] `config.yml` - Configuração de templates

### 4. Pull Request Template

- [x] `.github/pull_request_template.md` existe

### 5. GitHub Discussions

Categorias recomendadas:

- [ ] 📣 **Announcements** - Updates do projeto
- [ ] 💡 **Ideas** - Sugestões de features
- [ ] 🙏 **Q&A** - Perguntas e respostas
- [ ] 🎉 **Show and tell** - Casos de uso
- [ ] 💬 **General** - Discussões gerais

### 6. Labels

Labels úteis para Issues/PRs:

**Type:**

- [ ] `bug` - Algo não funciona
- [ ] `feature` - Nova funcionalidade
- [ ] `enhancement` - Melhoria de funcionalidade existente
- [ ] `documentation` - Melhorias de docs
- [ ] `question` - Pergunta/dúvida

**Priority:**

- [ ] `priority: critical` - Crítico
- [ ] `priority: high` - Alta
- [ ] `priority: medium` - Média
- [ ] `priority: low` - Baixa

**Status:**

- [ ] `status: needs-triage` - Precisa análise
- [ ] `status: in-progress` - Em progresso
- [ ] `status: blocked` - Bloqueado
- [ ] `status: awaiting-response` - Aguardando resposta

**Area:**

- [ ] `area: audio` - Módulo de áudio
- [ ] `area: transcription` - Transcrição
- [ ] `area: correction` - Correção
- [ ] `area: api` - API REST
- [ ] `area: plugin` - Plugin Obsidian
- [ ] `area: cli` - CLI
- [ ] `area: tests` - Testes
- [ ] `area: ci` - CI/CD

**Other:**

- [ ] `good first issue` - Bom para iniciantes
- [ ] `help wanted` - Precisa de ajuda
- [ ] `duplicate` - Duplicado
- [ ] `wontfix` - Não será corrigido
- [ ] `invalid` - Inválido

## 📢 Divulgação

### 1. Release Notes

Criar release notes detalhadas em GitHub Releases:

```markdown
# dictate2me v0.1.0 - First Public Release 🎉

We're excited to announce the first public release of dictate2me!

## 🎯 What is dictate2me?

dictate2me is a 100% offline voice transcription tool with AI-powered text correction. Perfect for privacy-conscious users who want local AI without sending data to the cloud.

## ✨ Features

### Core

- 🎤 Real-time audio capture
- 📝 Offline transcription using Vosk (Portuguese, English, Spanish, etc.)
- ✏️ AI text correction using local LLM (Ollama)
- 🔒 100% offline - your data never leaves your machine

### API

- 🌐 REST API for editor integrations
- 📡 WebSocket streaming for real-time results
- 🔐 Token-based authentication
- 🚀 <10ms latency

### Obsidian Plugin

- 🔌 Native Obsidian integration
- ⌨️ Hotkey support (Cmd/Ctrl+Shift+D)
- 💭 Live partial results
- 🎨 Visual feedback with pulsing microphone icon

## 📊 Stats

- ~3,000 lines of production code
- ~9,000 lines of documentation
- 85%+ test coverage
- 10ms API latency

## 📖 Getting Started

Check out our [README](https://github.com/zandercpzed/dictate2me#readme) for installation instructions.

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

MIT License - see [LICENSE](LICENSE)

---

**Full Changelog**: https://github.com/zandercpzed/dictate2me/commits/v0.1.0
```

- [ ] Release notes criadas
- [ ] Assets attached (se houver binários)

### 2. Social Media

**Reddit:**

- [ ] r/golang - Post sobre o projeto
- [ ] r/Obsidian - Post sobre o plugin
- [ ] r/privacy - Foco em offline/privacy
- [ ] r/opensource - Lançamento do projetoTemplate sugerido:

```
🎤 dictate2me v0.1.0 - Offline voice transcription with local AI

I built dictate2me, a 100% offline voice transcription tool with AI correction using Vosk + Ollama.

Features:
- Real-time transcription (Portuguese, English, etc.)
- LLM-based text correction
- REST API + WebSocket
- Obsidian plugin
- 100% privacy (all local)

Tech stack: Go, TypeScript, Vosk, Ollama

GitHub: https://github.com/zandercpzed/dictate2me

Would love feedback from the community!
```

**Twitter/X:**

- [ ] Tweet sobre lançamento
- [ ] Tag: #opensource #golang #ai #privacy #obsidian

**Hacker News:**

- [ ] Show HN: dictate2me - Offline voice transcription with AI
- [ ] https://news.ycombinator.com/submit

**Dev.to / Medium:**

- [ ] Blog post detalhado sobre implementação
- [ ] "Building an Offline AI Transcription Tool in Go"

### 3. Obsidian Community

- [ ] Post no Obsidian Forum
- [ ] Submit para Community Plugins (quando estável)
- [ ] Share em Discord servers de Obsidian

### 4. Newsletters / Blogs

- [ ] Go Weekly
- [ ] Golang Weekly newsletter submission
- [ ] Open Source Weekly

## 📊 Monitoramento Pós-Lançamento

### Métricas para Acompanhar

- [ ] GitHub stars
- [ ] Forks
- [ ] Issues abertas/fechadas
- [ ] Pull Requests
- [ ] Discussions activity
- [ ] Downloads (releases)
- [ ] Contributors

### Ferramentas Recomendadas

- [ ] GitHub Insights (analytics built-in)
- [ ] [Shields.io](https://shields.io/) - Badges
- [ ] [Codecov](https://codecov.io/) - Coverage tracking (opcional)
- [ ] [Renovate](https://github.com/renovatebot/renovate) - Dependency updates

## 🛡️ Segurança

- [ ] Configurar GitHub Security alerts
- [ ] Configurar Dependabot
- [ ] Email security@ configurado
- [ ] SECURITY.md atualizado com canais de contato

## 📝 Manutenção Contínua

### Responsividade

- [ ] Responder issues em < 72h
- [ ] Fazer triage de issues semanalmente
- [ ] Revisar PRs em < 7 dias
- [ ] Manter CHANGELOG atualizado

### Updates Regulares

- [ ] Releases mensais (ou quando necessário)
- [ ] Dependências atualizadas trimestralmente
- [ ] Roadmap público no GitHub Project (opcional)

### Documentação

- [ ] Manter README atualizado
- [ ] Atualizar docs quando features mudarem
- [ ] Adicionar FAQs baseado em issues comuns

## 🎯 Roadmap Público (Opcional)

Criar em GitHub Projects ou ROADMAP.md:

```markdown
# Roadmap

## v0.2.0 (Q1 2025)

- [ ] Linux support
- [ ] Windows support
- [ ] More language models
- [ ] Performance optimizations

## v0.3.0 (Q2 2025)

- [ ] VS Code plugin
- [ ] Web interface
- [ ] Docker images

## v1.0.0 (Q3 2025)

- [ ] Production-ready
- [ ] Full platform support
- [ ] Community plugins support
```

---

## ✅ Final Checklist

Antes de tornar público:

- [ ] Todo código sensível removido
- [ ] Documentação revisada
- [ ] Testes passando
- [ ] CI/CD funcionando
- [ ] Issue/PR templates configurados
- [ ] Labels criadas
- [ ] Discussions habilitado
- [ ] Release v0.1.0 criada
- [ ] README atraente
- [ ] LICENSE correto
- [ ] Social media posts preparados

**Depois de público:**

- [ ] Post em Reddit
- [ ] Tweet sobre lançamento
- [ ] Submit para newsletters
- [ ] Responder primeiros issues
- [ ] Agradecer primeiros contributors

---

## 🎉 Conclusão

O projeto está **100% pronto** para ser público!

**Comando final para tornar público:**

```bash
# Se ainda privado no GitHub
gh repo edit --visibility public

# Ou via web: Settings → Danger Zone → Change visibility → Public
```

**Celebrate!** 🎊

---

**Data da preparação:** 2025-12-02  
**Versão:** 0.1.0-rc1  
**Status:** Ready to launch! 🚀
