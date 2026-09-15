# AGENTS.md

Guidance for coding agents working in this repo.

## What this is

Provenance is the capture-enforcement gate at the code-to-push boundary. It verifies that the reasoning behind a change was actually captured in the vault before that change can be pushed, so capture stops being a good intention and becomes a precondition for shipping.

## Build and test

```bash
go build -o provenance ./cmd/provenance && go test ./...
```

## Conventions

- **The git hook is a convenience; the CI check is the gate.** A local hook can be skipped with `--no-verify`. If only the hook is installed, you have a nudge rather than a gate, and a nudge is ignorable exactly when capture matters most.
- **The handoff is a contract.** Canary consumes the declared provenance and vault governance for a change to rank risk and enforce never-skip overrides. Changing its shape is a cross-repo change.
- **Fail loudly at the boundary, never silently.** A gate that degrades quietly is worse than no gate, because it produces confidence without coverage.

## What not to do

- Do not add a bypass that is not visible in CI output. An invisible bypass is indistinguishable from a gate that never ran.
