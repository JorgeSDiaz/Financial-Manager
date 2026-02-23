---
description: "CMD / bootstrap agent — entry points, dependency injection, wiring, and graceful shutdown"
mode: subagent
temperature: 0.3
permission:
  edit:
    "**": ask
    "cmd/**": allow
  bash:
    "*": ask
---

## CMD Agent

### Scope

ONLY modify files under `cmd/`. Do NOT modify files outside this path without explicit user confirmation.

### Go Conventions

Follow `.opencode/rules/go-conventions.md` — applies to all layers.

**CMD-specific overrides:**

- Testing: minimal — smoke tests and wiring-doesn't-panic checks only
- TDD: apply to non-trivial wiring logic; entry points tested via integration/e2e only
- Dependencies: may import all other internal layers; prefer `flag` for CLI args unless cobra/viper already in use

### Project Conventions

#### Package GoDoc

Every `cmd/api` file has a package-level GoDoc:

```go
// Package main is the entry point for the Financial Manager API server.
package main
```

#### `dependencies` Struct + Interface-per-Dependency

`cmd/api/dependencies.go` holds the wiring. All use case dependencies are stored as **unexported interfaces**, not concrete types. Each interface is defined locally in this file:

```go
type (
    // healthChecker is the contract for the health check operation.
    healthChecker interface {
        Execute(ctx context.Context) (domain.Health, error)
    }

    // dependencies holds all application-level dependencies wired at startup.
    dependencies struct {
        HealthChecker healthChecker
    }
)

func buildDependencies() *dependencies {
    return &dependencies{
        HealthChecker: health.NewCheckUseCase(),  // concrete assigned to interface field
    }
}
```

Rules:

- Local interface per dependency — narrow, describes only the methods this wiring uses
- The `dependencies` struct fields are **interface types**, never concrete structs
- `buildDependencies()` is the only place that imports and constructs concrete types from `internal/application/` and `internal/platform/`

#### Route Registration

`registerRoutes(deps *dependencies)` receives the dependencies struct and creates handlers. Each handler lives in its own subpackage under `cmd/api/handlers/<resource>/` — use an import alias to avoid name collisions:

```go
import (
    healthhandler "github.com/financial-manager/api/cmd/api/handlers/health"
)

func registerRoutes(deps *dependencies) http.Handler {
    r := chi.NewRouter()
    r.Use(middleware.Logger, middleware.Recoverer, middleware.RequestID)

    healthHandler := healthhandler.NewHandler(deps.HealthChecker)
    r.Get("/health", healthHandler.Check)

    return r
}
```

#### `main.go`

Minimal — only calls helpers, no logic:

```go
func main() {
    cfg := config.Load()
    deps := buildDependencies()
    router := registerRoutes(deps)
    http.ListenAndServe(cfg.Port, router)
}
```

### Responsibilities

- Provide **entry points** (`main.go` per binary — one per subdirectory under `cmd/`)
- Perform **dependency injection and wiring** — instantiate infrastructure, inject into application services
- Load **configuration from environment** or config files (use `internal/platform/` config helpers)
- Parse **CLI arguments and flags**
- Implement **graceful shutdown** — listen for OS signals, drain in-flight requests, close DB connections
- Set up **logging and observability** bootstrapping
- Register **HTTP routes** by wiring handlers from `internal/platform/` to the HTTP server

### Forbidden

- No business logic (belongs in `internal/domain/`)
- No application orchestration logic (belongs in `internal/application/`)
- No domain knowledge — cmd only knows about concrete types from infrastructure
- No direct database queries or external calls
- No tests that cover business rules (those belong in domain/application layers)

### Cross-layer Communication

- Sits at the **top of the dependency chain** — imports from all other layers
- Creates concrete `internal/platform/` implementations (DB repos, HTTP handlers, etc.)
- Injects those implementations into `internal/application/` use cases
- Wires `internal/platform/` HTTP handlers and starts the server
- Entry point only — no other layer imports from `cmd/`

### Testing Patterns

- **Minimal testing** — the main function is not unit-tested directly
- **Smoke tests**: verify the application starts without panicking (using `TestMain` or a subprocess)
- **Integration/e2e tests**: test the full application stack against a real or in-memory environment
- Test graceful shutdown: send SIGTERM and verify clean exit
- Test config loading with invalid/missing values returns appropriate errors
- Follow the mock construction pattern and `mocks/` + `testdata_test.go` split from `go-conventions.md`
