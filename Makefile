.PHONY: build run dev test test-go test-web test-python lint lint-go lint-web lint-python fmt fmt-go fmt-web fmt-python generate check-generate clean frontend

# ── Build ──────────────────────────────────────────────

build: generate frontend
	@mkdir -p bin
	go build -o bin/autograder ./cmd/server

run: build
	./bin/autograder

dev:
	@echo "Starting Redis, Go server, and Vite dev server..."
	@echo "Open http://localhost:3000 for hot-reload dev"
	@echo "─────────────────────────────────────────────────"
	docker compose up -d
	@trap 'kill 0' INT; \
		go run ./cmd/server & \
		(cd web && deno task dev) & \
		wait

frontend:
	cd web && deno task build

# ── Code generation (FlatBuffers) ──────────────────────

generate:
	./schema/generate.sh

check-generate: generate
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "ERROR: 'make generate' produced unstaged changes. Run 'make generate' and commit."; \
		git status --short; \
		git diff --stat; \
		exit 1; \
	fi

# ── Test ───────────────────────────────────────────────

test: check-generate test-go test-web test-python

test-go:
	go test ./... -v -count=1

test-web:
	cd web && deno task test

test-python:
	cd marker && poetry install --quiet && poetry run pytest tests -v

# ── Lint (check only) ─────────────────────────────────

lint: lint-go lint-web lint-python

lint-go:
	@echo "── Go lint ──"
	go vet ./...
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './internal/models/*'))" || \
		(echo "gofmt needed on:"; gofmt -l $$(find . -name '*.go' -not -path './internal/models/*'); exit 1)

lint-web:
	@echo "── Web lint ──"
	cd web && deno task lint
	cd web && deno task fmt:check

lint-python:
	@echo "── Python lint ──"
	cd marker && poetry install --quiet
	cd marker && poetry run yapf --diff --recursive --exclude 'marker/generated/*' marker/ tests/
	cd marker && poetry run isort --check-only --diff marker/ tests/

# ── Format (auto-fix) ─────────────────────────────────

fmt: fmt-go fmt-web fmt-python

fmt-go:
	gofmt -w .

fmt-web:
	cd web && deno task lint:fix
	cd web && deno task fmt

fmt-python:
	cd marker && poetry install --quiet
	cd marker && poetry run yapf --in-place --recursive --exclude 'marker/generated/*' marker/ tests/
	cd marker && poetry run isort marker/ tests/

# ── Clean ──────────────────────────────────────────────

clean:
	rm -rf bin/ web/dist/ files/
