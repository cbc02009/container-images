# syntax=docker/dockerfile:1

# Paseo (https://paseo.sh) — a control room for coding agents. Extends the
# official daemon image with the OpenCode CLI and a mise-managed home-ops
# toolchain (kubectl, flux, kustomize, helm, sops, age, talosctl, topf, ...)
# so agent sessions can work on the GitOps repo directly.
#
# The daemon runs with its stock defaults (listen 0.0.0.0:6767, web UI
# enabled, no relay, no auth until PASEO_PASSWORD is set): it is meant to
# sit behind the internal gateway only.

ARG VERSION=0.15.4

FROM ghcr.io/getpaseo/paseo:${VERSION}@sha256:d9c05730c66a6b64a708303e7cfa6a4009f31f89507bf986fb5b098af25a1c3f

ARG VERSION
ARG OPENCODE_VERSION=2.3.0
# renovate: datasource=github-releases depName=jdx/mise
ARG MISE_VERSION=v2026.11.8
ARG TARGETARCH

USER root

RUN apt-get update && apt-get install -y --no-install-recommends \
        libatomic1 \
        ripgrep \
        unzip \
        xz-utils \
    && rm -rf /var/lib/apt/lists/*

RUN case "${TARGETARCH}" in \
        'amd64') OPENCODE_TARGET='linux-x64-baseline' ;; \
        'arm64') OPENCODE_TARGET='linux-arm64' ;; \
        *) echo "unsupported arch ${TARGETARCH}" && exit 1 ;; \
    esac \
    && curl -fsSL "https://registry.npmjs.org/@opencode/cli-${OPENCODE_TARGET}/-/cli-${OPENCODE_TARGET}-${OPENCODE_VERSION}.tgz" -o /tmp/opencode.tgz \
    && tar -xzOf /tmp/opencode.tgz package/bin/opencode > /usr/local/bin/opencode \
    && chmod 0755 /usr/local/bin/opencode \
    && rm -f /tmp/opencode.tgz

RUN case "${TARGETARCH}" in \
        'amd64') MISE_ARCH='linux-x64' ;; \
        'arm64') MISE_ARCH='linux-arm64' ;; \
        *) echo "unsupported arch ${TARGETARCH}" && exit 1 ;; \
    esac \
    && curl -fsSL "https://github.com/jdx/mise/releases/download/${MISE_VERSION}/mise-${MISE_VERSION}-${MISE_ARCH}.tar.xz" -o /tmp/mise.tar.xz \
    && tar -xJf /tmp/mise.tar.xz -C /tmp \
    && install /tmp/mise/bin/mise /usr/local/bin/mise \
    && rm -rf /tmp/mise /tmp/mise.tar.xz

ENV MISE_DATA_DIR=/home/paseo/.local/share/mise \
    MISE_CACHE_DIR=/home/paseo/.cache/mise \
    MISE_STATE_DIR=/home/paseo/.local/state/mise \
    PATH=/home/paseo/.local/share/mise/shims:/opt/mise/shims:/usr/local/bin:/usr/bin:/bin

COPY mise.toml /etc/mise/config.toml

RUN mkdir -p /opt/mise /tmp/mise-cache /tmp/mise-state \
    && cd /etc/mise \
    && MISE_DATA_DIR=/opt/mise MISE_CACHE_DIR=/tmp/mise-cache MISE_STATE_DIR=/tmp/mise-state mise trust \
    && MISE_DATA_DIR=/opt/mise MISE_CACHE_DIR=/tmp/mise-cache MISE_STATE_DIR=/tmp/mise-state mise install \
    && MISE_DATA_DIR=/opt/mise MISE_CACHE_DIR=/tmp/mise-cache MISE_STATE_DIR=/tmp/mise-state mise reshim \
    && opencode --version \
    && kubectl version --client \
    && flux version --client \
    && kustomize version \
    && helm version --short \
    && sops --version \
    && age --version \
    && talosctl version \
    && topf --version \
    && just --version \
    && rm -rf /tmp/mise-cache /tmp/mise-state

# The base image already runs the daemon as the paseo user (uid 1000) via
# tini + paseo-docker-entrypoint; keep that behavior.
USER 1000:1000
