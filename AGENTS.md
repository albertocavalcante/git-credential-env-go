# git-credential-env-go

git-credential-env-go

## Commands

```bash
just setup     # Install hooks and tools
just build     # Build
just test      # Run tests
just lint      # Run linter
just fmt       # Format code
just check     # Run all quality checks (what CI runs)
```

## Structure

```
.github/       -- CI/CD workflows
tools/         -- Dev tool dependencies
```

## Conventions

- Semantic commits: `type(scope): description` (max 72 chars)
- Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- Format with gofmt before committing
- Lint with golangci-lint -- all warnings are errors
- Tests must pass with race detector enabled
- SonarCloud for code coverage and quality gates (never Codecov)

## Do NOT

- Add dependencies to the main go.mod for dev-only tools (use tools.go.mod)
- Commit secrets or .env files
- Skip lefthook hooks (no --no-verify)
