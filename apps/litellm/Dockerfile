# syntax=docker/dockerfile:1

# LiteLLM overlay image: upstream base + three opencode_go provider fixes.
#
# 1. /v1/messages streaming: upstream's OpenCodeGoMessagesChatConfig does not
# override get_model_response_iterator_for_request, so streaming an OpenAI
# /v1/chat/completions wire request to a Go model whose endpoint is
# /v1/messages (e.g. opencode_go/minimax-m3) resolves to the pass-stub in
# BaseConfig, returns None as the stream iterator, and dies mid-stream with
# `TypeError: 'NoneType' object is not an iterator` → APIConnectionError 500
# (and, via the proxy router, a model failover). The COPY below adds the
# override, delegating to the Anthropic ModelResponseIterator — the same
# pattern OpenCodeGoGeminiChatConfig uses for the Gemini route. Tracked
# upstream at https://github.com/BerriAI/litellm/issues/51666.
#
# 2. Replayed reasoning fields: Go's OpenAI-compatible endpoint
# (/v1/chat/completions) strict-validates per-message fields and 400s with
# `unsupported_parameter: messages[N].reasoning` on any conversation whose
# assistant turns carry reasoning. litellm's drop_params only strips
# top-level request params — fields nested inside message objects pass
# through OpenAILikeChatConfig untouched, so the COPY below adds an
# OpenCodeGoChatConfig.transform_request override that strips
# reasoning/reasoning_content from the message list before forwarding.
#
# 3. Missing session header: Go answers 400 missing_session_id when a request
# carries no session id. litellm only attaches x-opencode-session when the
# caller established a session (x-litellm-session-id / metadata.session_id);
# there is no proxy-side config knob for headerless callers (the
# general_settings.missing_session_id "generate" policy is deliberately
# excluded from provider affinity). The COPY below wraps go_headers so all
# three Go wire configs synthesize a throwaway oneshot-<uuid> session id when
# the caller supplied none, leaving established sessions untouched.
#
# The opencode module is byte-identical across 1.109.0-rc.1, v1.109.0 and
# v1.109.1, so all three fixes are still required on this base; re-verify
# against the stock module before dropping either.

ARG VERSION=v1.109.1

FROM ghcr.io/berriai/litellm:${VERSION}@sha256:11175975a3093e05e92faee2fe2fe3412a3ae2fcaa849eb0fe25d5f0d8bcb2e9

# The venv path is python-version-pinned in the upstream image (3.13 at this
# tag); revisit when the base image changes its Python version.
COPY transformation.py /app/.venv/lib/python3.13/site-packages/litellm/llms/opencode/chat/transformation.py
