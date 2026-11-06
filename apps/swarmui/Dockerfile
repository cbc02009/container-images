# syntax=docker/dockerfile:1

# SwarmUI (https://github.com/mcmonkeyprojects/SwarmUI) — a self-hosted
# Stable Diffusion Web UI, built from source at a tagged release and served
# on port 7001.

# renovate: datasource=docker depName=mcr.microsoft.com/dotnet/sdk
FROM mcr.microsoft.com/dotnet/sdk:9.0-bookworm-slim

ARG UID=1000
ARG VERSION

RUN \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ffmpeg \
        git \
        libgl1 \
        libglib2.0-0 \
        python3.11 \
        python3.11-dev \
        python3.11-venv \
        wget \
    && rm -rf /var/lib/apt/lists/*

RUN \
    git config --global --add safe.directory '*' \
    && git clone --branch ${VERSION} --recurse-submodules --single-branch https://github.com/mcmonkeyprojects/SwarmUI.git /SwarmUI \
    && chown -R ${UID}:${UID} /SwarmUI

EXPOSE 7001

ENTRYPOINT ["bash", "/SwarmUI/launchtools/docker-standard-inner.sh", "fixperms"]

LABEL org.opencontainers.image.source="https://github.com/mcmonkeyprojects/SwarmUI"
