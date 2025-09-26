# Repository Guidelines

## Project Structure & Module Organization

- `cmd/server/main.go` bootstraps the HTTP server and loads configuration; keep additional entrypoints under `cmd/<name>/`.
- `internal/config` holds environment loading and database wiring; extend it for new providers without leaking implementation details elsewhere.
- Route wiring belongs in `internal/routers`; group HTTP handlers by resource under future `internal/handlers` subpackages.
- Persistence models and data access live in `internal/models` and `internal/repositories`; keep migrations and seeds alongside them for discoverability.
- Do not add unnessery comments

## Build, Test, and Development Commands

- `make build` builds containers via Docker Compose for a clean, reproducible stack.
- `make start` / `make stop` manage the running stack without rebuilding.
- `make logs` tails the application container stream; add `APP_PORT=4000 make start` to override defaults temporarily.
- `go run cmd/server` works for local, non-container hacking.

## Coding Style & Naming Conventions

- Rely on `gofmt` defaults (tabs, gofmt import ordering); run `golangci-lint fmt` before committing.
- Use UpperCamelCase for exported types/functions, lowerCamelCase for locals, and keep filenames lowercase with underscores only when separating words matters.
- Configuration keys mirror `.env` variables (e.g., `DB_USER` ⇒ `cfg.DBUser`).

## Testing Guidelines

- Prefer table-driven `*_test.go` files colocated with the code under `internal/`; reserve `tests/` for integration suites invoked by `make test`.
- Run `go test ./...` locally before pushing; aim to keep coverage meaningful around data access and router wiring.
- Mark external-service calls with interfaces so they can be faked in tests.

## Commit & Pull Request Guidelines

- Follow the existing imperative style (`init commit`); keep scope-focused summaries under ~60 characters.
- Reference tracking issues in the body using `Fixes #ID` when applicable and note environment or schema changes.
- Include screenshots or curl examples for new endpoints and list any new env vars in PR descriptions.

## Security & Configuration Tips

- Duplicate `.env.example` when onboarding; never commit real secrets.
- Rotate database credentials after sharing them and ensure health checks still succeed under new settings.
