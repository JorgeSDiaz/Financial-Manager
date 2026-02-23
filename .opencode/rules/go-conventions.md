# Go Conventions — Financial Manager API

Applies to all layers: `domain/`, `application/`, `platform/`, `cmd/`.

---

## Error Handling

- Wrap: `fmt.Errorf("operation: %w", err)`
- Inspect: `errors.Is()` / `errors.As()` only — never `==` on errors
- `errcheck ./...` must pass clean — no silently discarded errors

## Naming

- Exported identifiers: PascalCase. Unexported: camelCase
- Package names: lowercase, no underscores, no plurals (`health`, not `healths`)
- Interfaces named after behavior: `Reader`, `Fetcher`, `AccountRepository`
- Every exported type and function must have a GoDoc comment

## Imports

Three groups separated by blank lines, enforced by `goimports`:

```go
import (
    "context"      // 1. stdlib
    "fmt"

    "github.com/financial-manager/api/internal/domain" // 2. internal

    "github.com/go-chi/chi/v5" // 3. external
)
```

## `type()` Block Ordering

All type declarations in a single `type()` block per file.

**Domain files:**

1. Custom scalar types (`type HealthStatus string`)
2. Structs

**Handler / dependencies files:**

1. Interfaces (unexported, consumer contracts)
2. Handler or struct
3. DTOs (request/response, unexported)

Every entry in the block must have a GoDoc comment.

## `const()` Block

Group constants in a single `const()` block. Each constant must have a GoDoc comment.

## Constructors

- Named `NewXxx(...)`
- Domain entities/value objects → return **value type**
- Use cases, handlers → return **pointer**
- Always accept **interfaces**, never concrete types (platform/cmd layers)
- Validate invariants; return `error` on invalid state (domain layer)

## Testing

- **TDD**: Red → Green → Refactor
- Table-driven with anonymous structs and `t.Run()` subtests
- `testify/assert` (non-fatal) or `testify/require` (fatal)
- Test files: `*_test.go`. Fixtures: `testdata/` next to the test file
- Always pass `-race`. Minimum **90% coverage** per package
- Cover: happy path + edge cases + error paths

## Code Quality

- `gofmt -w .` before every commit
- `go vet ./...` — zero issues
- `golangci-lint run` — zero issues (`errcheck`, `gocyclo ≤10` included)
- No global variables — use dependency injection
- Prefer stdlib; minimal external dependencies; all imports from `go.mod`

## Pre-commit Checklist

```bash
gofmt -w .
go vet ./...
golangci-lint run
go test ./... -race
go mod tidy
```

---

## Directory Structure — One Subdirectory per Resource

Every layer organises code in one subdirectory per resource. The subdirectory name is the resource name in lowercase and matches the package name.

```
internal/domain/
└── health/
    └── health.go          # package health

internal/application/
└── health/
    ├── check.go           # package health
    └── check_test.go

cmd/api/handlers/
└── health/
    ├── handler.go         # package health
    ├── handler_test.go
    ├── testdata_test.go
    └── mocks/
        └── checker_mock.go
```

Rules:

- One package per resource — package name equals directory name (`package health`, not `package healthhandler`)
- Import aliases resolve name collisions at the call site: `healthhandler "github.com/financial-manager/api/cmd/api/handlers/health"`
- All other layers import the full path — no dot imports

---

## Interface-at-Consumer

**Interfaces are defined in the file that uses them, not in the package that implements them.**

- A handler defines its own narrow interface for each dependency (e.g. `checker`)
- The use case in `internal/application/` does **not** define its own interface
- This avoids import cycles, keeps interfaces minimal, and follows Go's implicit interface convention

```go
type (
    // checker is the contract for the health check operation.
    checker interface {
        Check(ctx context.Context) (domainhealth.Health, error)
    }

    // Handler serves HTTP health check requests.
    Handler struct {
        checker checker
    }
)
```

Constructor accepts the interface, never the concrete type:

```go
// NewHandler creates a new Handler with the given checker dependency.
func NewHandler(c checker) *Handler {
    return &Handler{checker: c}
}
```

---

## Test Infrastructure — `mocks/` and `testdata_test.go`

Test infrastructure within a resource package is split into two locations:

```
cmd/api/handlers/health/
├── handler.go            # implementation
├── handler_test.go       # TestXxx functions only
├── testdata_test.go      # build* helpers, fixture constructors, test doubles
└── mocks/
    └── checker_mock.go   # package mocks — mock type definitions only
```

**`mocks/<name>_mock.go`** (`package mocks`) — mock struct and interface method implementations only. No `*testing.T` dependency:

```go
// Checker is a mock implementation of the checker interface.
type Checker struct{ mock.Mock }

// Check mocks the checker.Check method.
func (m *Checker) Check(ctx context.Context) (domainhealth.Health, error) {
    args := m.Called(ctx)
    return args.Get(0).(domainhealth.Health), args.Error(1)
}
```

**`testdata_test.go`** (`package <resource>_test`) — everything that requires `*testing.T`: `build*` constructors, fixture value builders, test doubles:

```go
// buildMockChecker creates a Checker pre-configured to return result once.
func buildMockChecker(t *testing.T, result domainhealth.Health) *mocks.Checker {
    t.Helper()
    m := &mocks.Checker{}
    m.On("Check", mock.Anything).Return(result, nil).Once()
    return m
}
```

**Coverage exclusion** — `mocks/` must be excluded from coverage computation:

```makefile
go test ./... -coverprofile=coverage.out
grep -v '/mocks' coverage.out > coverage_filtered.out
go tool cover -func=coverage_filtered.out
```

The `make test` target handles this automatically.

---

## Mock Construction Pattern

`mock.On(...)` setup never lives inline in `t.Run`. The pattern has two parts:

**1. Builder function** — encapsulates construction and expectation setup:

```go
// buildMockChecker creates a Checker pre-configured to return result once.
func buildMockChecker(t *testing.T, result domainhealth.Health) *mocks.Checker {
    t.Helper()
    m := &mocks.Checker{}
    m.On("Check", mock.Anything).Return(result, nil).Once()
    return m
}
```

**2. Test table field** — each case declares its own `build<Mock>` func, invoked inside `t.Run`:

```go
tests := []struct {
    name         string
    buildChecker func(t *testing.T) *mocks.Checker
    wantStatus   int
}{
    {
        name: "returns 200 OK with StatusUp",
        buildChecker: func(t *testing.T) *mocks.Checker {
            return buildMockChecker(t, newHealthResult(domainhealth.StatusUp))
        },
        wantStatus: http.StatusOK,
    },
}

for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        checker := tc.buildChecker(t)
        // ... act and assert ...
        checker.AssertExpectations(t)
    })
}
```

Rules:

- Builder name: `build<Mock><Interface>(t *testing.T, ...)` — always `t.Helper()` as first call
- Encapsulates `On`, `Return`, and cardinality (`.Once()`, `.Times(n)`)
- Test case field type is `func(t *testing.T) *mocks.X` — never a pre-built instance
- `AssertExpectations(t)` always called after act, never omitted
