# syntax=docker/dockerfile:1

# Lens (https://github.com/BerriAI/lens) — the agent trace recording and
# investigation service that the LiteLLM dashboard's Lens sidebar connects to.
#
# Upstream publishes no release artifacts yet (no image, no chart release; see
# their deploy/lens/README.md), so this builds the standalone runtime from the
# pinned upstream tag by mirroring that tag's deploy/runtime/Dockerfile. When
# upstream changes its Dockerfile or starts publishing images/charts, re-sync
# this file — or drop the build in favor of the published, digest-verified
# artifact. The LENS_*_BUILD_IMAGE / LENS_*_RUNTIME_IMAGE arguments are
# digest-pinned upstream and are deliberately left WITHOUT renovate comments
# so Renovate does not fight upstream's pins; LENS_REF is the renovate-tracked
# upstream version.
#
# Upstream builds against the repository checkout as build context. This build
# cannot assume that, so a `lens-src` scratch stage is assembled from a clone
# of the pinned tag and every COPY of repository content is rewritten to
# `COPY --from=lens-src`. The smoke stage (sandbox test binary) from upstream
# is omitted: its runtime is exercised by container_test.go in PRs instead.

ARG LENS_REF=v1.108.0-rc.1
ARG LENS_BUILD_IMAGE=cgr.dev/chainguard/wolfi-base@sha256:1d95114038f76513a9ace6fca107d5582b08c65981f81f61cb56bf7fd2ef216d
ARG LENS_RUNTIME_IMAGE=cgr.dev/chainguard/wolfi-base@sha256:1d95114038f76513a9ace6fca107d5582b08c65981f81f61cb56bf7fd2ef216d
ARG LENS_UI_BUILD_IMAGE=node:24.14.1-bookworm-slim@sha256:b506e7321f176aae77317f99d67a24b272c1f09f1d10f1761f2773447d8da26c
ARG LENS_AGENT_BUILD_IMAGE=node:24.20.0-bookworm-slim@sha256:ba849c60be29959425b8734d57b8b4b7d56f98edd9504c9af091d5281095a71e

# renovate: datasource=docker depName=alpine/git
ARG GIT_IMAGE=alpine/git:v2.54.0@sha256:6a1fdb962814ca11508ca8d1bab94375b5b7ac4e4a9a40981c29fc32c02f095d

FROM ${GIT_IMAGE} AS source
ARG LENS_REF
RUN git clone --depth 1 --branch "${LENS_REF}" https://github.com/BerriAI/lens.git /src

FROM scratch AS lens-src
COPY --from=source /src /src

FROM ${LENS_UI_BUILD_IMAGE} AS ui-builder
WORKDIR /src
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=lens-src /src/package.json /src/package-lock.json ./
COPY --from=lens-src /src/src/ui/lib/package.json src/ui/lib/package.json
COPY --from=lens-src /src/src/ui/app/package.json src/ui/app/package.json
RUN npm ci
COPY --from=lens-src /src/src/ui/ src/ui/
RUN npm run build:ui

FROM ${LENS_AGENT_BUILD_IMAGE} AS agent-builder
WORKDIR /app/agent
COPY --from=lens-src /src/src/agent/package.json /src/src/agent/package-lock.json ./
RUN npm ci --ignore-scripts
COPY --from=lens-src /src/src/agent/ ./
RUN npm run build
WORKDIR /app/agent-chat
COPY --from=lens-src /src/src/agent-chat/package.json /src/src/agent-chat/package-lock.json ./
RUN npm ci --ignore-scripts
COPY --from=lens-src /src/src/agent-chat/ ./
RUN npm run build && npm prune --omit=dev --ignore-scripts && npm prune --prefix ../agent --omit=dev --ignore-scripts && \
    rm -f dist/*.test.js dist/*.test.d.ts dist/connectors/slack/*.test.js dist/connectors/slack/*.test.d.ts ../agent/dist/*.test.js ../agent/dist/*.test.d.ts

FROM ${LENS_BUILD_IMAGE} AS builder
RUN apk add --no-cache rust build-base cmake perl pkgconf openssl-dev libseccomp-dev python-3.13 jq git
WORKDIR /src
COPY --from=lens-src /src/.cargo/ .cargo/
COPY --from=lens-src /src/src/worker/ src/worker/
WORKDIR /src/src/worker
ENV CARGO_BUILD_JOBS=2 CARGO_PROFILE_RELEASE_DEBUG=0 CARGO_PROFILE_RELEASE_STRIP=symbols \
    CARGO_PROFILE_RELEASE_LTO=thin CARGO_PROFILE_RELEASE_CODEGEN_UNITS=16
RUN cargo build --locked --release -p litellm-lens
RUN cargo build --locked --release -p lens-migrate
RUN cc -std=c11 -D_GNU_SOURCE -O2 -Wall -Wextra -Werror crates/lens/sandbox/python_policy.c -lseccomp -o /tmp/python-policy && \
    /tmp/python-policy /tmp/python.seccomp

FROM ${LENS_RUNTIME_IMAGE} AS service
ARG LENS_VERSION
RUN apk add --no-cache python-3.13 setpriv libgcc libstdc++ openssl ca-certificates curl
ENV LENS_VERSION=${LENS_VERSION} PYTHONDONTWRITEBYTECODE=1 LENS_UI_DIRECTORY=/app/ui
WORKDIR /app
COPY --from=builder /src/src/worker/target/release/litellm-lens /usr/local/bin/litellm-lens
COPY --from=builder /src/src/worker/target/release/lens-migrate /usr/local/bin/lens-migrate
COPY --from=builder /tmp/python.seccomp /app/lens/python.seccomp
COPY --from=ui-builder /src/src/ui/app/out/ /app/ui/
COPY --from=agent-builder /usr/local/bin/node /usr/local/bin/node
COPY --from=agent-builder /app/agent/ /app/agent/
COPY --from=agent-builder /app/agent-chat/ /app/agent-chat/
COPY --from=lens-src /src/src/worker/crates/lens/sandbox/python_runtime.py /tmp/python_runtime.py
RUN python3.13 -S /tmp/python_runtime.py /app/lens/python-runtime.json && rm /tmp/python_runtime.py
USER 65532:65532
EXPOSE 4318
HEALTHCHECK --interval=10s --timeout=3s --start-period=30s --retries=6 CMD curl --fail --silent http://127.0.0.1:4318/health/ready > /dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/node", "/app/agent-chat/dist/supervisor.js"]

FROM service AS runtime
