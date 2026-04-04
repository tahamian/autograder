FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=C.UTF-8

# ── System packages ────────────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    git \
    unzip \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# ── FlatBuffers (latest from GitHub) ───────────────────
ARG FLATC_VERSION=25.2.10
RUN curl -fsSL "https://github.com/google/flatbuffers/releases/download/v${FLATC_VERSION}/Linux.flatc.binary.g++-13.zip" \
    -o /tmp/flatc.zip \
    && unzip /tmp/flatc.zip -d /usr/local/bin/ \
    && chmod +x /usr/local/bin/flatc \
    && rm /tmp/flatc.zip

# ── Go ─────────────────────────────────────────────────
ARG GO_VERSION=1.26.1
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    | tar -xz -C /usr/local
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# ── Deno ───────────────────────────────────────────────
RUN curl -fsSL https://deno.land/install.sh | sh
ENV DENO_INSTALL="/root/.deno"
ENV PATH="${DENO_INSTALL}/bin:${PATH}"

# ── Python + Poetry ────────────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-venv \
    && rm -rf /var/lib/apt/lists/*
RUN pip3 install --break-system-packages poetry

# ── Verify ─────────────────────────────────────────────
RUN go version && deno --version && python3 --version && poetry --version && flatc --version

WORKDIR /workspace
