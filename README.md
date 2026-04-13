# qube

> Declarative API testing CLI.
> One command to test HTTP, gRPC, GraphQL, WebSocket and more.

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![CI](https://github.com/apiqube/qube/actions/workflows/ci.yml/badge.svg)](https://github.com/apiqube/qube/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen?style=flat-square)]()

## Install

```bash
go install github.com/apiqube/qube/cmd@latest
```

## Quick Start

```yaml
# tests/hello.yaml
target: http://localhost:8081

tests:
  - name: Health check
    method: GET
    endpoint: /health
    expect:
      status: 200
```

```bash
qube run tests/
```

## Commands

```
qube run [path]       Run tests
qube check [path]     Validate without executing
qube init             Create starter files
qube generate         Generate tests from OpenAPI/Swagger
qube plugin           Manage plugins
qube version          Print version
```

## Related Repositories

| Repo | Description |
|---|---|
| [`apiqube/engine`](https://github.com/apiqube/engine) | Core engine library |
| [`apiqube/cli`](https://github.com/apiqube/cli) | V1 CLI (archived reference) |

## License

[MIT](LICENSE)
