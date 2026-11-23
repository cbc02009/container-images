# syntax=docker/dockerfile:1

# Official TrueNAS MCP server (https://github.com/truenas/truenas-mcp) —
# the upstream release workflow publishes binaries only (no container image),
# so this clones the pinned release tag, builds the official binary the same
# way upstream's release does (CGO_ENABLED=0, version injected via
# -X main.Version), and ships it on distroless static. Upstream embeds IANA
# tzdata (_ "time/tzdata") for NAS-timezone next_run, so the runtime needs no
# tzdata package; distroless static supplies the CA bundle TLS verification
# loads, and the server intentionally exits at startup unless it can log in
# to a TrueNAS running 25.10 or later — it cannot serve without a NAS.

ARG VERSION=1.0.1

# renovate: datasource=docker depName=alpine/git
ARG TRUENAS_MCP_GIT_IMAGE=alpine/git:v2.54.0@sha256:6a1fdb962814ca11508ca8d1bab94375b5b7ac4e4a9a40981c29fc32c02f095d
# renovate: datasource=docker depName=golang:1.27.4-alpine3.23
ARG TRUENAS_MCP_BUILD_IMAGE=golang:1.27.4-alpine3.23@sha256:f70850ef6ce341e8b6427f677049145d63c0b7ce002d4d7e3d464828b1b96dcd
# renovate: datasource=docker depName=gcr.io/distroless/static-debian12
ARG TRUENAS_MCP_RUNTIME_IMAGE=gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

FROM ${TRUENAS_MCP_GIT_IMAGE} AS source
ARG VERSION
RUN git clone --depth 1 --branch "v${VERSION}" https://github.com/truenas/truenas-mcp.git /src

FROM scratch AS truenas-src
COPY --from=source /src /src

FROM ${TRUENAS_MCP_BUILD_IMAGE} AS build
WORKDIR /src

# Dependencies first, so source edits do not invalidate the module cache.
COPY --from=truenas-src /src/go.mod /src/go.sum ./
RUN go mod download

COPY --from=truenas-src /src/ .

ARG VERSION
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" \
    -o /out/truenas-mcp ./cmd/truenas-mcp

FROM ${TRUENAS_MCP_RUNTIME_IMAGE}

COPY --from=build /out/truenas-mcp /usr/local/bin/truenas-mcp

USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/truenas-mcp"]
