# Testes Unitários — Dictate2Me

## Como rodar os testes

Execute no terminal, dentro da pasta do projeto:

```bash
npx ts-node src/audio.test.ts
npx ts-node src/transcription.test.ts
npx ts-node src/correction.test.ts
```

## O que validar

- Logs estruturados de cada módulo
- Integração entre captura, transcrição e correção
- Saída esperada:
    - `[TEST][AUDIO] { dataCalled: true }`
    - `[TEST][TRANSCRIPTION] { result: ... }`
    - `[TEST][CORRECTION] { result: ... }`

---

## Teste de Integração (MVP)

1. Instale e ative o plugin no Obsidian.
2. Abra uma nota qualquer.
3. Clique no ícone de microfone na ribbon ou use o comando "Iniciar transcrição de voz".
4. Observe o status bar mudando para "Gravando..." e depois exibindo o texto corrigido.
5. Verifique se o texto foi inserido no cursor do editor ativo.
6. Confira os logs estruturados no console do Obsidian (View → Toggle Developer Tools).
7. Repita o teste usando o comando da paleta.

## Próximos passos

- Adicionar testes de integração automatizados (futuro)
- Cobertura de comandos e fluxo completo
