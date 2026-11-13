# syntax=docker/dockerfile:1

# LiteLLM overlay image: upstream base + the opencode_go /v1/messages
# streaming fix. Upstream's OpenCodeGoMessagesChatConfig does not override
# get_model_response_iterator_for_request, so streaming an OpenAI
# /v1/chat/completions wire request to a Go model whose endpoint is
# /v1/messages (e.g. opencode_go/minimax-m3) resolves to the pass-stub in
# BaseConfig, returns None as the stream iterator, and dies mid-stream with
# `TypeError: 'NoneType' object is not an iterator` → APIConnectionError 500
# (and, via the proxy router, a model failover). The COPY below adds the
# override, delegating to the Anthropic ModelResponseIterator — the same
# pattern OpenCodeGoGeminiChatConfig uses for the Gemini route. Tracked
# upstream at https://github.com/BerriAI/litellm/issues/51666 — drop this
# overlay once a release ships that fix.

ARG VERSION=v1.109.0-rc.1

FROM ghcr.io/berriai/litellm:${VERSION}@sha256:030121401c60d732a2cdf8fd485c91432d40388b046bfa3e0b2d3370f26378c9

# The venv path is python-version-pinned in the upstream image (3.13 at this
# tag); revisit when the base image changes its Python version.
COPY transformation.py /app/.venv/lib/python3.13/site-packages/litellm/llms/opencode/chat/transformation.py
