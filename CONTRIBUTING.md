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

## 🎉 Reconhecimento

Todos os contribuidores serão reconhecidos no arquivo [CONTRIBUTORS.md](CONTRIBUTORS.md).

---

Dúvidas? Abra uma [Discussion](https://github.com/zandercpzed/dictate2me/discussions)!
