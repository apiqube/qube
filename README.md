# qube

> Declarative API testing CLI.
> One command to test HTTP, gRPC, GraphQL, WebSocket and more.

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![CI](https://github.com/apiqube/qube/actions/workflows/ci.yml/badge.svg)](https://github.com/apiqube/qube/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen?style=flat-square)]()

`qube` is the official CLI for [ApiQube](https://github.com/apiqube). It runs declarative test manifests through the engine, drives a live Bubble Tea UI when you have a terminal, and produces JSON / JUnit / TAP for CI.

## Install

```bash
go install github.com/apiqube/qube/cmd@latest
```

## Quick start

```bash
qube init                # scaffold .qube.yaml + tests/example.yaml
qube run tests/          # execute the suite
qube check tests/        # validate without running
qube plugin list         # what protocols are available
```

```yaml
# tests/example.yaml
target: http://localhost:8080

tests:
  - name: Health check
    method: GET
    resource: /health
    expect:
      status: 200
```

## Commands

| Command | Status | What it does |
|---|---|---|
| `qube run [path]` | ✅ live | Run tests with a live TUI; `--output=json|junit|tap` for machine-readable output |
| `qube check [path]` | ✅ live | Validate manifests without executing |
| `qube init` | ✅ live | Create starter files; `--interactive` launches a wizard |
| `qube plugin list` | ✅ live | Show installed WASM plugins |
| `qube version` | ✅ live | Print build info |
| `qube generate` | 🚧 stub | Generate tests from OpenAPI/Swagger/HAR/Postman (roadmap) |
| `qube plugin install` | 🚧 stub | Install plugins from registry (roadmap) |
| `qube plugin remove` | 🚧 stub | Remove installed plugins (roadmap) |

## Output modes

`qube run --output` selects how results are reported:

- `pretty` (default) — Bubble Tea live UI on a terminal; auto-fallback to progressive lipgloss output when piped or in CI
- `json` — newline-delimited JSON, one event per line
- `junit` — JUnit XML on `RunCompleted` (GitHub Actions / GitLab CI ready)
- `tap` — Test Anything Protocol with YAML diagnostics

## Configuration

`.qube.yaml` is auto-discovered by walking up from the current directory:

```yaml
version: 1
targets:
  default: http://localhost:8080
runner:
  parallel: true
  failFast: false
plugins:
  - http
```

`.env` files in the same hierarchy populate `{{ env.* }}` template references in tests.

## Plugins

`qube` drives WASM plugins for protocol support. The first-party `plugin-http` ships in [`apiqube/plugin-http`](https://github.com/apiqube/plugin-http). Drop a `.wasm` file in `~/.apiqube/plugins/` (or set `$QUBE_PLUGIN_DIR`) and `qube plugin list` shows it.

## Related repositories

| Repo | Description |
|---|---|
| [`apiqube/engine`](https://github.com/apiqube/engine) | Core engine library (declarative testing runtime) |
| [`apiqube/plugin-http`](https://github.com/apiqube/plugin-http) | First-party HTTP plugin |
| [`apiqube/cli`](https://github.com/apiqube/cli) | V1 CLI (archived reference) |

## License

[MIT](LICENSE)
