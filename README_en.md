<div align="center">
  <img src="https://avatars.githubusercontent.com/u/245985800?s=200&v=4" alt="logo" width="100" />
  <h2>Example</h2>
  <h3>Complete FasterEdge Usage Example (Go Demo)</h3>
</div>

### 1. Introduction

- **Purpose**: A **complete runnable example** built on the [FasterEdge](https://github.com/FasterEdge/FasterEdge) main project. It demonstrates initializing an Atom, mounting built-in abilities, invoking ability commands, executing the authenticated-command path, and managing the component lifecycle.
- **Dependency**: `github.com/FasterEdge/FasterEdge v1.0.20261022` (see `go.mod`)

### 2. Quick Start

```bash
# 1. Fetch dependencies (including the FasterEdge main project)
go mod download

# 2. Run the demo
go run .

# 3. Build a native executable
go build .
```

### 3. What the Demo Shows

| Ability | Command | Description |
|---|---|---|
| **TimeAbility** | `get_time` | Fetch the current time snapshot (source / time) |
| **NetMapAbility** | `register_peer` / `list_peers` | Register a peer and query the network topology |
| **OneKeyAbility** | `issue_token` / `verify_token` | Issue a short-lived access token for a peer and verify it |
| **CmdAbility** | `set_allowlist` / `run` | Configure a command-execution allowlist and run a controlled command |
| **Auth path** | `AuthenticatedCommand` | Invoke commands through the security boundary with a OneKey credential; tampered signatures are rejected |
| **Lifecycle** | `RunAtom` / `CloseAtom` | Run all Runner abilities, then shut down gracefully on cancellation |

### 4. Cross-Platform Auto-Build

The repository ships a GitHub Actions workflow (`.github/workflows/release.yml`):

- Trigger: push a `v*` tag or manual `workflow_dispatch`
- Build matrix: **Windows / macOS / Linux × amd64 / arm64** (6 targets)
- Release: packages `tar.gz` (Unix) and `zip` (Windows) and publishes them to a GitHub Release

```bash
# Tag to trigger the automated build and release
git tag v1.0.0 && git push origin v1.0.0
```

### 5. License

- **License**: The root `LICENSE` file contains the Apache License 2.0.
- **Application**: Example source retains the FasterEdge organization attribution and copyright notice.
