# Guia de Contribuição

Obrigado por considerar contribuir com o dictate2me! Este documento fornece diretrizes para contribuir com o projeto.

## 📋 Índice

- [Código de Conduta](#código-de-conduta)
- [Como Posso Contribuir?](#como-posso-contribuir)
- [Configurando o Ambiente](#configurando-o-ambiente)
- [Padrões de Código](#padrões-de-código)
- [Processo de Pull Request](#processo-de-pull-request)
- [Conventional Commits](#conventional-commits)

## 📜 Código de Conduta

Este projeto adota o [Código de Conduta do Contributor Covenant](CODE_OF_CONDUCT.md). Ao participar, espera-se que você mantenha este código.

## 🤔 Como Posso Contribuir?

### Reportando Bugs

Antes de criar um bug report:

1. Verifique se o bug já não foi reportado em [Issues](https://github.com/zandercpzed/dictate2me/issues)
2. Se não encontrar, crie uma issue usando o template de bug report

### Sugerindo Melhorias

Sugestões são sempre bem-vindas! Use o template de feature request.

### Contribuindo com Código

1. Procure issues marcadas com `good first issue` ou `help wanted`
2. Comente na issue que você gostaria de trabalhar nela
3. Aguarde um mantenedor atribuir a issue a você

### Melhorando a Documentação

Documentação é tão importante quanto código. PRs de documentação são muito valorizados.

## 🛠️ Configurando o Ambiente

### Pré-requisitos

- Go 1.23+
- Git
- Make ou Mage
- golangci-lint
- pre-commit

### Setup

```bash
# Clone o repositório
git clone https://github.com/zandercpzed/dictate2me.git
cd dictate2me

# Execute o script de setup
./scripts/setup-dev.sh

# Instale os hooks de pré-commit
pre-commit install

# Verifique se tudo está funcionando
make test
```

## 📝 Padrões de Código

### Go

- Siga o [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` para formatação
- Todas as funções públicas DEVEM ter comentários GoDoc
- Cobertura de testes: 100% é obrigatório

### Comentários

```go
// TranscribeAudio transcreve um arquivo de áudio para texto.
//
// O arquivo deve estar no formato WAV, 16kHz, mono, 16-bit.
// Retorna o texto transcrito e um erro se a transcrição falhar.
//
// Exemplo:
//
//	text, err := TranscribeAudio("audio.wav")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(text)
func TranscribeAudio(path string) (string, error) {
    // implementação
}
```

### Testes

```go
func TestTranscribeAudio(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "arquivo válido em português",
            input:   "testdata/audio/sample-pt-br.wav",
            want:    "olá mundo",
            wantErr: false,
        },
        {
            name:    "arquivo inexistente",
            input:   "nonexistent.wav",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := TranscribeAudio(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("TranscribeAudio() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("TranscribeAudio() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 🔄 Processo de Pull Request

1. **Fork** o repositório
2. **Clone** seu fork localmente
3. **Crie uma branch** para sua feature/fix:
   ```bash
   git checkout -b feat/minha-feature
   ```
4. **Faça commits** seguindo Conventional Commits
5. **Execute os testes** localmente:
   ```bash
   make test
   make lint
   ```
6. **Push** para seu fork
7. **Abra um PR** para a branch `main`

### Checklist do PR

- [ ] Código segue os padrões do projeto
- [ ] Testes adicionados/atualizados
- [ ] Cobertura de testes mantida em 100%
- [ ] Documentação atualizada
- [ ] Commits seguem Conventional Commits
- [ ] PR tem descrição clara do que foi feito

## 📌 Conventional Commits

Usamos [Conventional Commits](https://www.conventionalcommits.org/) para mensagens de commit.

### Formato

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Tipos

| Tipo       | Descrição                         |
| ---------- | --------------------------------- |
| `feat`     | Nova funcionalidade               |
| `fix`      | Correção de bug                   |
| `docs`     | Apenas documentação               |
| `style`    | Formatação, sem mudança de código |
| `refactor` | Refatoração de código             |
| `perf`     | Melhoria de performance           |
| `test`     | Adição ou correção de testes      |
| `build`    | Mudanças no build system          |
| `ci`       | Mudanças no CI                    |
| `chore`    | Outras mudanças                   |

### Exemplos

```bash
feat(audio): add voice activity detection
fix(transcription): handle empty audio files gracefully
docs: update installation instructions for macOS
test(correction): add tests for Portuguese grammar rules
```

## 🏗️ Arquitetura e Design

### Princípios Arquiteturais

1. **Modularidade**: Cada módulo deve ser independente e testável
2. **Interfaces Claras**: Use interfaces Go para abstrair dependências
3. **Simplicidade**: Prefira soluções simples a complexas
4. **Performance**: Código deve ser eficiente, mas legível
5. **Offline-First**: Tudo deve funcionar 100% offline

### Estrutura de Diretórios

```
dictate2me/
├── cmd/                    # Binários executáveis
│   ├── dictate2me/        # CLI principal
│   └── dictate2me-daemon/ # Daemon da API
├── internal/              # Código interno (não importável)
│   ├── audio/            # Captura de áudio
│   ├── transcription/    # Motor de transcrição
│   ├── correction/       # Correção de texto
│   └── api/              # API REST
├── pkg/                   # Código público (importável)
├── docs/                  # Documentação
├── scripts/               # Scripts utilitários
├── test/                  # Testes de integração
└── plugins/               # Plugins (Obsidian, etc.)
```

### ADRs (Architecture Decision Records)

Decisões arquiteturais importantes devem ser documentadas em `docs/adr/`.

**Template:** Use `docs/adr/template.md`

**Quando criar um ADR:**

- Mudança de tecnologia (ex: trocar Whisper por Vosk)
- Nova funcionalidade significativa (ex: adicionar WebSocket)
- Mudança de arquitetura (ex: adicionar cache)

## 🔍 Code Review Process

### Para Revisores

**Checklist:**

- [ ] Código segue style guide
- [ ] Testes cobrem casos principais
- [ ] Documentação está atualizada
- [ ] Sem breaking changes (ou bem documentados)
- [ ] Performance não foi degradada
- [ ] Segurança não foi comprometida

**Feedback:**

- Seja construtivo e educado
- Explique o "porquê" das sugestões
- Aponte o que está bom também
- Sugira melhorias, não exija

### Para Contribuidores

**Respondendo ao review:**

- Agradeça o feedback
- Faça perguntas se não entender
- Implemente ou discuta sugestões
- Marque conversas como resolvidas

## 📖 Documentação

### Tipos de Documentação

1. **Code Comments (GoDoc)**

   - Todas as funções/tipos públicos
   - Explique o "o quê" e "porquê"
   - Inclua exemplos quando útil

2. **README.md**

   - Para cada submódulo
   - Quickstart e exemplos
   - Atualizar quando mudar comportamento

3. **ADRs (`docs/adr/`)**

   - Decisões arquiteturais
   - Contexto, decisão, consequências

4. **Guides (`docs/`)**
   - Tutoriais passo-a-passo
   - Troubleshooting
   - Architecture overview

### Style Guide de Documentação

````markdown
# Título Principal (H1)

Breve descrição em 1-2 sentenças.

## Seção (H2)

### Subseção (H3)

**Negrito** para destaque.
_Itálico_ para ênfase.
`código inline` para comandos/código.

```bash
# Blocos de código com syntax highlighting
```
````

- Listas com `-`
- Não use `*` ou `+`

1. Listas numeradas
2. Quando ordem importa

````

## 🧪 Testes - Guia Detalhado

### Cobertura Obrigatória

- **Novos pacotes**: 90%+ coverage
- **Funções críticas**: 100% coverage
- **Edge cases**: Sempre teste

### Estratégia de Testes

```go
// 1. Table-Driven Tests (preferido)
func TestMyFunc(t *testing.T) {
    tests := []struct{
        name string
        // ...
    }{
        // casos de teste
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}

// 2. Subtests
func TestComplex(t *testing.T) {
    t.Run("happy path", func(t *testing.T) { /*...*/ })
    t.Run("error case", func(t *testing.T) { /*...*/ })
}

// 3. Setup/Teardown
func TestWithSetup(t *testing.T) {
    setup := createTestFixture()
    defer setup.Cleanup()
    // ...
}
````

### Mocking

```go
// Use interface para dependências
type Transcriber interface {
    Transcribe([]int16) (string, error)
}

// Mock em teste
type MockTranscriber struct {
    TranscribeFunc func([]int16) (string, error)
}

func (m *MockTranscriber) Transcribe(audio []int16) (string, error) {
    return m.TranscribeFunc(audio)
}
```

## 🚀 Release Process

### Versioning

Usamos [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH`
- MAJOR: Breaking changes
- MINOR: Novas features (backward compatible)
- PATCH: Bug fixes

### Processo de Release

1. **Atualizar CHANGELOG.md**

   ```markdown
   ## [1.2.0] - 2025-MM-DD

   ### Added

   - Nova feature X

   ### Fixed

   - Bug Y
   ```

2. **Criar tag**

   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```

3. **CI cria release automaticamente**

   - Build de binários
   - Publicação no GitHub Releases

4. **Anunciar**
   - GitHub Discussions
   - Se maior: blog post

## 🐛 Debugging

### Logs

```go
// Use slog para logging estruturado
import "log/slog"

slog.Info("transcription started",
    "model", modelName,
    "duration", duration)

slog.Error("transcription failed",
    "error", err,
    "audio_size", len(audio))
```

### Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace
go test -trace=trace.out
go tool trace trace.out
```

## 🔒 Security

### Reporting Vulnerabilities

**NUNCA** reporte vulnerabilidades via issue pública.

Use: security@dictate2me.dev

### Security Checklist

PR com código sensível deve garantir:

- [ ] Sem hardcoded secrets
- [ ] Input validation adequada
- [ ] Sem SQL injection (se aplicável)
- [ ] Dependências atualizadas
- [ ] Sem logs de dados sensíveis

## 💬 Comunicação

### Onde Discutir

| Tópico      | Canal                   |
| ----------- | ----------------------- |
| Bugs        | GitHub Issues           |
| Features    | GitHub Discussions      |
| Arquitetura | GitHub Discussions      |
| Dúvidas     | GitHub Discussions Q&A  |
| Segurança   | security@dictate2me.dev |

### Etiqueta

- 📝 Seja claro e conciso
- 🤝 Seja respeitoso
- 🔍 Pesquise antes de perguntar
- 💡 Compartilhe conhecimento
- 🎉 Celebre sucessos da comunidade

## 🎓 Recursos para Aprender

### Go

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

### Testes em Go

- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [Testify](https://github.com/stretchr/testify)

### Contribuindo Open Source

- [First Timers Only](https://www.firsttimersonly.com/)
- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)

## 🎉 Reconhecimento

Todos os contribuidores serão reconhecidos:

- ✅ [CONTRIBUTORS.md](CONTRIBUTORS.md) - Lista de todos
- ✅ Release notes - Créditos por feature
- ✅ GitHub Contributors graph

### Tipos de Contribuição

Reconhecemos TODAS as formas de contribuição:

- 💻 Código
- 📖 Documentação
- 🐛 Bug reports
- 💡 Ideas
- 🎨 Design
- 🌍 Traduções
- 🧪 Testes
- 📣 Divulgação

---

## 📝 FAQs

**Q: Quanto tempo leva para um PR ser revisado?**  
A: Geralmente 1-3 dias úteis. PRs maiores podem levar mais tempo.

**Q: Posso trabalhar em múltiplas issues ao mesmo tempo?**  
A: Recomendamos focar em uma de cada vez para facilitar o review.

**Q: O que fazer se meu PR ficar desatualizado?**  
A: Faça rebase ou merge da branch main e resolva conflitos.

**Q: Posso contribuir se sou iniciante em Go?**  
A: Sim! Procure issues marcadas com `good first issue`.

---

Dúvidas? Abra uma [Discussion](https://github.com/zandercpzed/dictate2me/discussions)!
