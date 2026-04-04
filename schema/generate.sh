#!/usr/bin/env bash
# Generates Go, TypeScript, and Python models from the FlatBuffers schema.
# Run from the project root: ./schema/generate.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "==> Generating Go models..."
mkdir -p "$ROOT_DIR/internal/models"
flatc -g \
  --gen-object-api --gen-onefile \
  --go-namespace models \
  --go-module-name autograder/internal \
  -o "$ROOT_DIR/internal/models/" \
  "$ROOT_DIR/schema/models.fbs"

echo "==> Generating TypeScript models..."
mkdir -p "$ROOT_DIR/web/src/generated"
flatc -T \
  --gen-object-api --gen-all \
  -o "$ROOT_DIR/web/src/generated/" \
  "$ROOT_DIR/schema/models.fbs"

echo "==> Generating Python models..."
mkdir -p "$ROOT_DIR/marker/marker/generated"
flatc -p \
  --gen-object-api --gen-onefile \
  --python-typing \
  -o "$ROOT_DIR/marker/marker/generated/" \
  "$ROOT_DIR/schema/models.fbs"
touch "$ROOT_DIR/marker/marker/generated/__init__.py"

echo "==> Done."
