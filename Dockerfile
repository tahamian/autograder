# ── Stage 1: Build frontend ────────────────────────────
FROM denoland/deno:2.7.11 AS frontend

WORKDIR /app/web
COPY web/ .
RUN deno task build

# ── Stage 2: Build Go binary ──────────────────────────
FROM golang:1.26-alpine AS backend

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY schema/ schema/
COPY --from=frontend /app/web/dist web/dist/

RUN CGO_ENABLED=0 go build -o /autograder ./cmd/server

# ── Stage 3: Runtime ──────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates docker-cli

WORKDIR /app
COPY --from=backend /autograder .
COPY --from=frontend /app/web/dist web/dist/
COPY config.yml .
COPY marker/ marker/

EXPOSE 9090

ENTRYPOINT ["./autograder"]
