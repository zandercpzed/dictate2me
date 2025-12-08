# Guia Técnico — Plugin Obsidian de Transcrição Local

## Decisões Técnicas

- Linguagem: TypeScript
- Sem módulos nativos
- Logs estruturados: `console.log('[MÓDULO]', { dados })`
- Build: esbuild
- Estrutura: src/main.ts como entrypoint

## Estrutura Recomendada

- package.json
- tsconfig.json
- esbuild.config.mjs
- manifest.json
- src/main.ts

## Padrões de Código

- Modularização por funcionalidade
- Tipos explícitos
- Sem dependências nativas
- Testes unitários (futuro)
