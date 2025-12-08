> **Privacy-first voice writing assistant for Obsidian**

---

## Visão do Produto

dictate2me é um plugin para Obsidian que transforma voz em texto de forma inteligente, processando **100% localmente** — sem dependência de serviços em nuvem, sem transmissão de dados, sem assinaturas mensais.

**Tagline:** _Sua voz, suas notas, seu controle._

---

## O Problema

Escritores, pesquisadores e profissionais do conhecimento enfrentam um dilema:

| Soluções em Nuvem               | Soluções Locais Existentes           |
| ------------------------------- | ------------------------------------ |
| Excelente qualidade             | Qualidade inferior                   |
| Zero privacidade                | Privacidade total                    |
| Dependência de internet         | Funcionam offline                    |
| Custos recorrentes ($10-20/mês) | Geralmente gratuitas                 |
| Dados sensíveis expostos        | Transcrição "crua" sem processamento |

**O gap:** Não existe uma solução que combine qualidade de nível comercial com privacidade absoluta e processamento inteligente de texto.

---

## Proposta de Valor

```
┌─────────────────────────────────────────────────────────┐
│                                                         │
│   Transcrição de alta qualidade                         │
│   + Processamento inteligente (pontuação, parágrafos)   │
│   + Privacidade total (100% local)                      │
│   + Custo zero após setup                               │
│   + Integração nativa com Obsidian                      │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## Público-Alvo

### Primário

- **Usuários avançados de Obsidian** que valorizam privacidade
- **Profissionais que lidam com informações sensíveis** (advogados, médicos, pesquisadores)
- **Escritores e criadores de conteúdo** que preferem ditar a digitar

### Secundário

- Desenvolvedores e técnicos que preferem soluções self-hosted
- Usuários em regiões com internet instável ou cara
- Pessoas com limitações físicas que dependem de ditado

### Persona Principal

> **Marina, 34, Pesquisadora Acadêmica** Usa Obsidian para Zettelkasten. Trabalha com dados de entrevistas confidenciais. Já testou WisprFlow mas abandonou por preocupações com privacidade. Tem um MacBook M1 com capacidade de processamento local.

---

## Funcionalidades Core

### MVP (v1.0)

| Funcionalidade                | Descrição                                      | Prioridade |
| ----------------------------- | ---------------------------------------------- | ---------- |
| **Transcrição Local**         | Speech-to-text via Whisper (Transformers.js)   | P0         |
| **Processamento Inteligente** | LLM local (Ollama) para pontuação e formatação | P0         |
| **Limpeza de Fala**           | Remoção automática de "éh", "hum", repetições  | P0         |
| **Detecção de Parágrafos**    | Segmentação inteligente por pausas e contexto  | P0         |
| **Hotkey Global**             | Atalho para iniciar/parar gravação             | P0         |
| **Inserção em Cursor**        | Texto processado inserido na posição atual     | P0         |

### Futuro (v2.0+)

- Comandos de voz para formatação ("novo parágrafo", "título")
- Vocabulário customizado por vault
- Múltiplos idiomas simultâneos
- Modo "append" para journaling contínuo
- Templates de processamento personalizáveis

---

## Arquitetura Técnica

```
┌─────────────────────────────────────────────────────────────────┐
│                        OBSIDIAN                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    DICTATE2ME PLUGIN                       │  │
│  │                                                           │  │
│  │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐   │  │
│  │   │   Audio     │───▶│   Whisper   │───▶│   Ollama    │   │  │
│  │   │   Capture   │    │ (local STT) │    │ (local LLM) │   │  │
│  │   └─────────────┘    └─────────────┘    └─────────────┘   │  │
│  │         │                   │                  │          │  │
│  │         ▼                   ▼                  ▼          │  │
│  │   ┌─────────────────────────────────────────────────┐     │  │
│  │   │              Pipeline de Processamento          │     │  │
│  │   │                                                 │     │  │
│  │   │  1. Captura de áudio (Web Audio API)           │     │  │
│  │   │  2. Transcrição bruta (Whisper via TFJS)       │     │  │
│  │   │  3. Processamento de texto (Ollama)            │     │  │
│  │   │     • Pontuação e capitalização                │     │  │
│  │   │     • Remoção de disfluências                  │     │  │
│  │   │     • Segmentação em parágrafos                │     │  │
│  │   │  4. Inserção no editor                         │     │  │
│  │   └─────────────────────────────────────────────────┘     │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │  TUDO PERMANECE   │
                    │     LOCAL         │
                    │  Zero transmissão │
                    └───────────────────┘
```

---

## Stack Tecnológico

| Camada             | Tecnologia                | Justificativa                                |
| ------------------ | ------------------------- | -------------------------------------------- |
| **Linguagem**      | TypeScript                | Padrão Obsidian, type-safety                 |
| **Speech-to-Text** | Whisper (Transformers.js) | Melhor modelo open-source, roda no browser   |
| **LLM Local**      | Ollama                    | Simplicidade de setup, ampla compatibilidade |
| **Modelo LLM**     | Llama 3.2 3B / Phi-3 Mini | Balanço qualidade/performance                |
| **Áudio**          | Web Audio API             | Nativo, sem dependências                     |
| **Estado**         | Obsidian Plugin API       | Integração nativa                            |

### Requisitos de Sistema

```
Mínimo:
├── RAM: 8GB
├── Storage: 5GB (modelos)
└── CPU: Qualquer processador moderno

Recomendado:
├── RAM: 16GB
├── GPU: Apple Silicon / NVIDIA com CUDA
└── CPU: Multi-core recente
```

---

## Fluxo do Usuário

```
    ┌──────────┐
    │  Usuário │
    │ pressiona│
    │  hotkey  │
    └────┬─────┘
         │
         ▼
    ┌──────────┐     ┌──────────┐     ┌──────────┐
    │ Gravação │────▶│Transcrição────▶│Processamento
    │  inicia  │     │  Whisper │     │  Ollama  │
    └──────────┘     └──────────┘     └────┬─────┘
                                           │
         ┌─────────────────────────────────┘
         │
         ▼
    ┌──────────┐     ┌──────────┐
    │  Texto   │────▶│ Inserido │
    │  limpo   │     │no cursor │
    └──────────┘     └──────────┘

    Latência alvo: < 3s após fim da fala
```

---

## Processamento de Texto — Especificação

### Entrada (transcrição bruta)

```
"éh então eu estava pensando que que a gente poderia fazer
um um sistema que funciona localmente sabe porque porque
a privacidade é muito importante"
```

### Saída (texto processado)

```
Então, eu estava pensando que a gente poderia fazer um
sistema que funciona localmente, sabe? Porque a privacidade
é muito importante.
```

### Regras de Processamento

1. **Disfluências removidas:** "éh", "hum", "tipo assim", repetições imediatas
2. **Pontuação inferida:** vírgulas, pontos, interrogações baseadas em entonação e contexto
3. **Capitalização:** início de frases, nomes próprios
4. **Parágrafos:** quebra após pausas longas (>2s) ou mudança de tópico
5. **Preservação:** gírias intencionais, estilo pessoal, estrutura argumentativa

---

## Diferenciais Competitivos

| Aspecto                   | WisprFlow | Apple Dictation | dictate2me    |
| ------------------------- | --------- | --------------- | ------------- |
| Privacidade               | ❌ Cloud  | ⚠️ Parcial      | ✅ 100% Local |
| Qualidade                 | ✅ Alta   | ✅ Alta         | ✅ Alta       |
| Custo                     | $10/mês   | Grátis          | Grátis        |
| Offline                   | ❌        | ⚠️ Limitado     | ✅ Total      |
| Processamento inteligente | ✅        | ❌              | ✅            |
| Integração Obsidian       | ❌        | ❌              | ✅ Nativa     |
| Customização              | ❌        | ❌              | ✅ Total      |

---

## Modelo de Distribuição

```
┌─────────────────────────────────────────┐
│            Open Source (MIT)            │
│                                         │
│  • Plugin gratuito via Community Plugins│
│  • Código aberto no GitHub              │
│  • Documentação completa                │
│  • Comunidade para suporte              │
│                                         │
└─────────────────────────────────────────┘
```

**Monetização potencial (futuro):**

- Modelos fine-tuned premium para domínios específicos
- Suporte empresarial
- Versão com UI simplificada para não-técnicos

---

## Métricas de Sucesso

### Técnicas

| Métrica               | Alvo MVP | Alvo v2.0 |
| --------------------- | -------- | --------- |
| Word Error Rate (WER) | < 10%    | < 5%      |
| Latência end-to-end   | < 5s     | < 2s      |
| Uso de RAM            | < 4GB    | < 2GB     |

### Produto

| Métrica                  | Alvo 6 meses | Alvo 12 meses |
| ------------------------ | ------------ | ------------- |
| Downloads                | 1.000        | 10.000        |
| GitHub Stars             | 500          | 2.000         |
| Usuários ativos semanais | 200          | 2.000         |
| Issues resolvidos        | 90% em 7d    | 95% em 3d     |

---

## Riscos e Mitigações

| Risco                                       | Probabilidade | Impacto | Mitigação                                                   |
| ------------------------------------------- | ------------- | ------- | ----------------------------------------------------------- |
| Performance insuficiente em hardware antigo | Alta          | Alto    | Múltiplos modelos Whisper (tiny→large), detecção automática |
| Setup complexo do Ollama                    | Média         | Alto    | Wizard de instalação, documentação visual                   |
| Qualidade inferior a cloud                  | Média         | Alto    | Fine-tuning, prompts otimizados, feedback loop              |
| Conflitos com outros plugins                | Baixa         | Médio   | Testes extensivos, API bem isolada                          |

---

## Roadmap

```
Q1 2025 ──────────────────────────────────────────────────────────
│
├── v0.1 Alpha
│   └── Transcrição básica funcionando
│
├── v0.5 Beta
│   └── Processamento com Ollama integrado
│
└── v1.0 Release
    └── MVP completo no Community Plugins

Q2 2025 ──────────────────────────────────────────────────────────
│
├── v1.1
│   └── Comandos de voz básicos
│
└── v1.2
    └── Multi-idioma

Q3 2025 ──────────────────────────────────────────────────────────
│
└── v2.0
    └── Vocabulário custom, templates, UX refinada
```

---

## Decisões Arquiteturais Chave

### ADR-001: Whisper via Transformers.js vs. API local

**Decisão:** Transformers.js  
**Motivo:** Elimina dependência externa, funciona em qualquer OS, setup mais simples

### ADR-002: Ollama vs. llama.cpp direto

**Decisão:** Ollama  
**Motivo:** API HTTP padrão, gerenciamento de modelos, comunidade ativa

### ADR-003: Processamento síncrono vs. streaming

**Decisão:** Híbrido — transcrição em chunks, processamento após fala completa  
**Motivo:** Balanço entre feedback visual e qualidade de processamento

---

## Chamada para Ação

> **Próximo passo:** Implementar POC da pipeline de transcrição com Whisper tiny model para validar performance em diferentes hardwares.

---

_Documento versão 1.0 — Dezembro 2024_  
_Projeto: dictate2me_
