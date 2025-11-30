# 📁 Project Structure Overview

This document provides a visual overview of the dictate2me project structure.

## Root Level

```
dictate2me/
├── 📄 README.md                  # Main documentation
├── 📄 LICENSE                    # MIT License
├── 📄 CONTRIBUTING.md            # Contribution guidelines
├── 📄 CODE_OF_CONDUCT.md        # Code of conduct
├── 📄 SECURITY.md               # Security policy
├── 📄 CHANGELOG.md              # Version history
├── 📄 GOVERNANCE.md             # Project governance
├── 📄 MAINTAINERS.md            # Maintainer list
├── 📄 SUPPORT.md                # Support information
├── 📄 BOOTSTRAP_COMPLETE.md     # Bootstrap phase summary
│
├── 📄 go.mod                    # Go module definition
├── 📄 Makefile                  # Build automation
├── 📄 .gitignore               # Git ignore rules
├── 📄 .editorconfig            # Editor configuration
├── 📄 .golangci.yaml           # Linter configuration
│
├── 📁 .github/                 # GitHub specific files
│   ├── 📁 ISSUE_TEMPLATE/      # Issue templates
│   ├── 📁 workflows/           # CI/CD workflows
│   └── 📄 PULL_REQUEST_TEMPLATE.md
│
├── 📁 cmd/                     # Application entry points
│   ├── 📁 dictate2me/          # CLI application
│   │   └── 📄 main.go
│   └── 📁 dictate2me-daemon/   # Daemon process
│       └── 📄 main.go
│
├── 📁 internal/                # Private application code
│   ├── 📁 audio/               # Audio capture
│   │   └── 📄 doc.go
│   ├── 📁 transcription/       # Speech-to-text
│   │   └── 📄 doc.go
│   ├── 📁 correction/          # Text correction
│   │   └── 📄 doc.go
│   ├── 📁 integration/         # Editor integrations
│   │   └── 📁 obsidian/
│   ├── 📁 api/                 # REST API
│   ├── 📁 config/              # Configuration
│   └── 📁 platform/            # OS-specific code
│
├── 📁 pkg/                     # Public reusable packages
│   └── 📁 textutils/
│
├── 📁 plugins/                 # Editor plugins
│   └── 📁 obsidian-dictate2me/
│       └── 📁 src/
│
├── 📁 models/                  # AI models (gitignored)
│   ├── 📄 README.md
│   └── 📄 .gitkeep
│
├── 📁 docs/                    # Documentation
│   ├── 📁 adr/                 # Architecture Decision Records
│   │   ├── 📄 README.md
│   │   ├── 📄 template.md
│   │   └── 📄 0001-linguagem-go.md
│   ├── 📁 blueprints/
│   ├── 📁 diagrams/
│   └── 📁 api/
│
├── 📁 scripts/                 # Utility scripts
│   └── 📄 setup-dev.sh         # Development setup
│
├── 📁 configs/                 # Configuration examples
│
├── 📁 testdata/                # Test fixtures
│   ├── 📁 audio/
│   │   └── 📄 .gitkeep
│   └── 📁 text/
│       └── 📄 .gitkeep
│
└── 📁 bin/                     # Built binaries (gitignored)
    ├── dictate2me
    └── dictate2me-daemon
```

## Key Directories

### `/cmd` - Application Entry Points

Contains the `main.go` files for executable applications. Each subdirectory represents a separate binary.

### `/internal` - Private Application Code

Code that is specific to this application and should not be imported by other projects.

- `audio/` - Audio capture and processing
- `transcription/` - Whisper integration
- `correction/` - LLM-based text correction
- `integration/` - Integrations with external tools
- `api/` - REST API server
- `config/` - Configuration management
- `platform/` - OS-specific implementations

### `/pkg` - Public Libraries

Code that can be imported by other projects. Keep this minimal.

### `/plugins` - Editor Plugins

Integrations with text editors and IDEs.

### `/models` - AI Models

Downloaded AI models (Whisper, LLM). Gitignored due to size.

### `/docs` - Documentation

- `adr/` - Architecture Decision Records
- `blueprints/` - Module design documents
- `diagrams/` - Visual diagrams
- `api/` - API documentation

### `/scripts` - Utility Scripts

Helper scripts for development, deployment, and maintenance.

### `/configs` - Configuration Examples

Example configuration files for users.

### `/testdata` - Test Fixtures

Test data used by test suites.

## File Naming Conventions

- `main.go` - Entry point for executables
- `doc.go` - Package documentation
- `*_test.go` - Test files
- `*_darwin.go` - macOS-specific code
- `*_linux.go` - Linux-specific code
- `*_windows.go` - Windows-specific code

## Import Paths

```go
import (
    // Internal packages
    "github.com/zandercpzed/dictate2me/internal/audio"
    "github.com/zandercpzed/dictate2me/internal/transcription"

    // Public packages (if any)
    "github.com/zandercpzed/dictate2me/pkg/textutils"
)
```

## Build Outputs

```
bin/
├── dictate2me           # CLI binary
└── dictate2me-daemon    # Daemon binary
```

---

**Last updated**: 2025-01-30  
**Version**: 0.0.1-bootstrap
