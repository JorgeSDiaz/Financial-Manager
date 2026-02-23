# Architecture — Financial Manager API

## Style

Hexagonal (Ports & Adapters).

## Layers

```
cmd/api/          → entry point, DI wiring, routes
internal/
  domain/         → entities, value objects, ports (repository interfaces)
  application/    → use cases
  platform/       → HTTP handlers, repos, adapters, config
```

## Dependency Rule

`cmd → platform → application → domain`. **No layer may import upward.**

- `domain/` imports nothing internal
- `application/` imports only `domain/`
- `platform/` imports `application/` and `domain/`
- `cmd/` imports all internal layers — it is the only wiring point

Violating this rule is a hard error. If you find yourself needing to import upward, the solution is to extract an interface (port) in the lower layer and inject the implementation from above.
