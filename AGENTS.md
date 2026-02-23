# AGENTS.md — Financial Manager API

> Go conventions: `.opencode/rules/go-conventions.md`  
> Architecture + dependency rule: `.opencode/rules/architecture.md`  
> Layer-specific rules: `.opencode/agents/`  
> Global quality standards: `~/.config/opencode/rules/always/code-quality.md`

Module: `github.com/financial-manager/api` | Runtime: Go 1.23 | Router: `github.com/go-chi/chi/v5`

---

## Commands

```bash
make run          # go run ./cmd/api
make build        # go build -o bin/api ./cmd/api
make test         # go test ./... -v -race -coverprofile=coverage.out
make lint         # golangci-lint run ./...
make tidy         # go mod tidy
```

**Single test:**

```bash
go test ./internal/domain/... -run TestNewHealth -v
go test ./internal/application/health/... -v -race
go test ./cmd/api/handlers/... -v -race
```

**Coverage:**

```bash
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```
