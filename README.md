# Autograder

[![CI](https://github.com/tahamian/autograder/actions/workflows/ci.yml/badge.svg)](https://github.com/tahamian/autograder/actions/workflows/ci.yml)

Autograder evaluates Python scripts uploaded via a web interface and grades them
against configurable test cases. Student code runs inside sandboxed Docker
containers with no network access, memory limits, and CPU constraints.

## Quick Start

Just Docker required — nothing else to install:

```bash
docker compose up
```

Open **http://localhost:9090** — select a lab, upload a `.py` file, and get graded.

This builds and starts everything automatically:

- **autograder** — Go API server + TypeScript frontend
- **marker** — Python grading image (installed as a pip package)
- **redis** — rate limiting

## How It Works

1. Student uploads a `.py` file and selects a lab
2. The Go server writes the file + test config to a submission directory
3. A `autograder-marker` Docker container runs the student code in a sandbox
4. The marker evaluates stdout and function return values against expected output
5. Results are returned as JSON and rendered in the frontend

## Architecture

```
┌─────────────┐     POST /api/submit     ┌──────────────┐
│   Browser    │ ──────────────────────── │   Go Server  │
│  (Vite/TS)  │ ◄─────── JSON ────────── │   :9090      │
└─────────────┘                          └──────┬───────┘
                                                │
                                    ┌───────────▼───────────┐
                                    │  autograder-marker     │
                                    │  Docker container      │
                                    │  (no network, 256MB)   │
                                    │                        │
                                    │  marker --config-file  │
                                    │    /mnt/input.json     │
                                    │    --output-file       │
                                    │    /mnt/output.json    │
                                    └────────────────────────┘
```

## Project Structure

```
autograder/
├── cmd/server/main.go               # Entry point
├── internal/
│   ├── api/                          # JSON API handlers, Redis rate limiter, SPA server
│   ├── config/                       # YAML config + env var overrides
│   ├── docker/                       # Docker client interface, image management, container lifecycle
│   ├── grader/                       # Evaluation logic (compares marker output to expected)
│   └── models/                       # Generated Go types (FlatBuffers)
├── schema/
│   ├── models.fbs                    # FlatBuffers schema (source of truth for shared types)
│   └── generate.sh                   # Generates Go, TypeScript, and Python types
├── web/                              # Frontend (Vite + TypeScript + Tailwind CSS, managed by Deno)
│   ├── deno.json                     # Deno config, tasks, dependencies
│   ├── vite.config.ts
│   ├── eslint.config.js              # ESLint + typescript-eslint
│   ├── .prettierrc.json              # Prettier config
│   ├── index.html                    # SPA shell (Tailwind classes)
│   └── src/
│       ├── models.ts                 # TypeScript interfaces matching FlatBuffers schema
│       ├── api.ts                    # API client (fetchLabs, submitFile)
│       ├── ui.ts                     # Testable DOM helpers (buildLabOptions, buildEvalCard, etc.)
│       ├── main.ts                   # App wiring
│       ├── style.css                 # Tailwind entry point
│       ├── generated/                # Generated TypeScript types (FlatBuffers)
│       ├── api.test.ts               # API tests (fetch mocking)
│       ├── ui.test.ts                # UI tests (jsdom + chai)
│       └── test-setup.ts             # jsdom DOM bootstrap
├── marker/                           # Python grading package (runs inside Docker)
│   ├── Dockerfile                    # Installs marker as a pip package
│   ├── pyproject.toml                # Poetry config + yapf/isort/pytest settings
│   ├── marker/
│   │   ├── __init__.py               # Package exports
│   │   ├── cli.py                    # Click CLI entry point (`marker` command)
│   │   ├── grader.py                 # Reads input JSON, runs Assignment, writes output
│   │   ├── sandbox.py                # Code execution sandbox (imports, functions, stdout)
│   │   └── generated/                # Generated Python types (FlatBuffers)
│   └── tests/
├── config.yml                        # Lab definitions + server settings
├── docker-compose.yml                # One command to run everything
├── Dockerfile                        # Multi-stage: Deno frontend → Go binary → Alpine runtime
├── Makefile                          # Build, test, lint, format, generate
├── go.mod
└── .github/
    ├── Dockerfile                    # CI image (Go + Deno + Python + Poetry + flatc)
    └── workflows/ci.yml
```

## Configuration

### Labs

Labs are defined in `config.yml`. Two test types are supported:

**stdout** — checks what the program prints:

```yaml
labs:
  - id: lab_1
    name: Hello World
    problem_statement: Write a program that prints "Hello World"
    testcase:
      - type: stdout
        name: hello_world
        expected:
          - feedback: "Correct!"
            points: 1.0
            values: ["Hello World", "hello world"]
```

**function** — imports and calls a function, checks the return value:

```yaml
labs:
  - id: lab_2
    name: Pythagorean Theorem
    problem_statement: Write pythagorean(a, b) returning the hypotenuse
    testcase:
      - type: function
        name: test_pythagorean
        functions:
          - function_name: pythagorean
            function_args:
              - { value: 3.0, type: float }
              - { value: 4.0, type: float }
        expected:
          - feedback: "Correct!"
            points: 1.0
            values: [5.0]
```

### Environment Variables

| Variable                    | Default                 | Description                              |
|-----------------------------|-------------------------|------------------------------------------|
| `REDIS_URL`                 | `redis://0.0.0.0:6379`  | Redis connection URL                     |
| `REDIS_RATE_LIMIT`          | `50-H`                  | Rate limit (requests per period)         |
| `REDIS_MAX_RETRY`           | `3`                     | Redis max retry attempts                 |
| `AUTOGRADER_HOST`           | `0.0.0.0`               | Server bind address                      |
| `AUTOGRADER_PORT`           | `9090`                  | Server port                              |
| `AUTOGRADER_MARKER_IMAGE`   | from `config.yml`       | Docker image for grading containers      |
| `AUTOGRADER_HOST_FILES_DIR` | —                       | Host path for Docker bind mounts         |

### API

| Method | Endpoint       | Description                              |
|--------|----------------|------------------------------------------|
| GET    | `/api/labs`    | List available labs                      |
| POST   | `/api/submit`  | Submit a file (`file` + `lab_id`)        |
| GET    | `/`            | Serves the frontend SPA                  |

## Development

### Prerequisites (for local dev without Docker)

| Tool         | Version | Install              |
|--------------|---------|----------------------|
| Go           | 1.26+   | `brew install go`    |
| Deno         | 2.x     | `brew install deno`  |
| Docker       | —       | docker.com           |
| FlatBuffers  | 25+     | `brew install flatbuffers` |
| Poetry       | 2.x     | `pip install poetry` |

### Development with hot-reload (Docker)

```bash
make dev
```

This starts everything in Docker with automatic rebuilds:

- **Go** — `air` watches `.go`/`.yml` files and rebuilds/restarts on change
- **Frontend** — Vite dev server with HMR on **http://localhost:3000** (proxies `/api` → `:9090`)
- **Marker** — Docker Compose `watch` rebuilds the image when `marker/` files change
- **Redis** — standard service

### Development without Docker (local tools)

```bash
make dev-local
```

Starts Redis in Docker, Go server and Vite dev server locally.
Requires all prerequisites installed on your machine.

### Production-like local run

```bash
# Build marker image
docker compose build marker

# Build frontend + backend and run
make run
```

### Code Generation (FlatBuffers)

Shared model types are defined once in `schema/models.fbs` and generated for
Go (`internal/models/`), TypeScript (`web/src/generated/`), and Python
(`marker/marker/generated/`):

```bash
make generate         # regenerate all
make check-generate   # verify no unstaged changes (used in CI)
```

### Makefile Targets

| Command            | Description                                              |
|--------------------|----------------------------------------------------------|
| `make run`         | Build everything and start the server                    |
| `make dev`         | Docker dev env with hot-reload (air + Vite HMR + watch)  |
| `make dev-local`   | Local dev with Redis in Docker, Go + Vite on host        |
| `make build`       | Generate models, build frontend, compile Go binary       |
| `make test`        | Check generation + Go + frontend + Python tests          |
| `make lint`        | Check all linting (gofmt, go vet, ESLint, Prettier, yapf, isort) |
| `make fmt`         | Auto-fix all formatting                                  |
| `make generate`    | Regenerate types from `schema/models.fbs`                |
| `make clean`       | Remove `bin/`, `web/dist/`, `files/`                     |

### Linting & Formatting

| Language   | Linter           | Formatter | Config                   |
|------------|------------------|-----------|--------------------------|
| Go         | `go vet`         | `gofmt`   | —                        |
| TypeScript | ESLint           | Prettier  | `eslint.config.js`, `.prettierrc.json` |
| Python     | isort            | yapf      | `pyproject.toml`         |

Generated code is excluded from all linting and formatting.

## Testing

**58 tests** across three languages:

| Component   | Tests | Framework                          |
|-------------|-------|------------------------------------|
| Go backend  | 32    | `go test` + mock Docker client     |
| TypeScript  | 22    | Deno test + chai + jsdom           |
| Python      | 4     | Poetry + pytest                    |

CI also checks FlatBuffers generation is up to date (no unstaged changes after
running `generate.sh`).

## CI

GitHub Actions with a custom Docker image (`.github/Dockerfile`) containing all
tools: Go, Deno, Python, Poetry, and flatc. Single `ci` job runs everything:
generation check → Go lint/test → frontend lint/format/test/build → Python
lint/test.

## License

MIT
