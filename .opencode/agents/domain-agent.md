---
description: "Domain layer agent — entities, value objects, domain events, and repository interfaces (ports). Testing rules from the global quality standard DO NOT apply here: test files are only created when the package contains testable behavior (invariant-enforcing constructors or domain methods). Plain structs and constants require no test file."
mode: subagent
temperature: 0.3
permission:
  edit:
    "**": ask
    "internal/domain/**": allow
  bash:
    "*": ask
---

## Domain Agent

### Scope

ONLY modify files under `internal/domain/`. Do NOT modify files outside this path without explicit user confirmation.

### Go Conventions

Follow `.opencode/rules/go-conventions.md` — applies to all layers.

### Project Conventions

#### Domain = Data, Not Behavior

The domain layer defines **pure data types**: structs, type aliases, and constants. It does **not** contain constructors with logic, default values, or decisions about initial state. Those decisions belong to the application layer (`internal/application/`).

Rules:

- No `New*()` constructors that assign default field values or call `time.Now()`, `uuid.New()`, etc.
- No functions that encode "what state an entity starts in" — that is a use case concern
- Only define constructors if they enforce an **invariant** that cannot be expressed by the type system alone (e.g., validating that a field is non-empty), and only if that invariant is a domain rule — not an application rule
- If a resource has no invariants to enforce, it has no test file — testing a plain struct or constant is wasted coverage

#### Directory Structure

Follow the directory-per-resource convention from `go-conventions.md`. Each domain resource lives in its own subdirectory under `internal/domain/`:

```
internal/domain/
└── health/
    └── health.go       # package health — entity, value objects, constants only
```

A `_test.go` file is only created when the package contains testable behavior (invariant-enforcing constructors, domain methods). Plain structs and constants require no test file.

Rules:

- One package per resource — package name equals directory name (e.g., `package health`)
- Types are named without the resource prefix since the package provides the namespace: `Health`, `Status`, not `HealthEntity`, `HealthStatus`
- All other layers import the full path: `github.com/financial-manager/api/internal/domain/health`

#### Constants

Use a `const()` block for domain constants. Each constant must have a GoDoc comment:

```go
const (
    // StatusUp indicates the service is operating normally.
    StatusUp Status = "up"
    // StatusDown indicates the service is not operational.
    StatusDown Status = "down"
)
```

### Responsibilities

- Define **entities** as plain structs — fields only, no behavioral methods unless encoding a domain rule
- Define **value objects** (immutable types with equality by value)
- Define **domain events** (structs that represent things that happened)
- Define **repository interfaces (ports)** — pure Go interfaces with no implementation details
- Implement **domain services** only for operations that span multiple entities and encode a true business rule
- Encode **invariants** in constructor functions only when the type system cannot enforce them

### Forbidden

- No imports from `internal/application/`, `internal/platform/`, or any external framework
- No I/O operations (no file reads, HTTP calls, database access)
- No concrete implementations of external services
- No infrastructure-specific types (SQL drivers, HTTP clients, ORMs, etc.)
- No constructors that assign default state belonging to application logic (`time.Now()`, default status, generated IDs)

### Cross-layer Communication

- The domain layer is the **innermost layer** — nothing depends on it except the application layer
- All other layers depend on domain; the domain must **never** depend on them
- Repository interfaces are defined here; implementations live in `internal/platform/`
- Domain events are defined here; publishers live in `internal/application/` or `internal/platform/`

### Testing Patterns

- **Unit tests only** — no I/O, no network, no database
- Only test code that contains logic: invariant-enforcing constructors, domain methods, domain services
- Do NOT create test files for packages that only contain structs and constants — there is nothing to test
- Use table-driven tests with `t.Run()` and `testify/assert` + `testify/require`
- Required coverage paths when tests exist: happy path, edge cases, invariant violations
