package main

import (
	"testing"

	helpers "github.com/cbc02009/container-images/tests"
)

// The image is an overlay on the upstream LiteLLM image carrying three
// opencode_go fixes: (1) OpenCodeGoMessagesChatConfig must define
// get_model_response_iterator_for_request itself (the inherited BaseConfig
// stub returns None, which kills OpenAI-wire streaming to /v1/messages Go
// models with 'NoneType' object is not an iterator); (2) OpenCodeGoChatConfig
// must strip per-message reasoning fields, which Go's /v1/chat/completions
// endpoint rejects with 400 unsupported_parameter; (3) every Go wire config
// must guarantee an x-opencode-session header, which Go otherwise rejects
// with 400 missing_session_id when the caller supplied no session id.
func TestLitellm(t *testing.T) {
	image := helpers.GetTestImage("litellm:v1.109.1")

	// The patched provider module replaced the stock one in the venv.
	helpers.RequireFileExists(
		t,
		image,
		"/app/.venv/lib/python3.13/site-packages/litellm/llms/opencode/chat/transformation.py",
	)

	// The override is defined on the Messages config class itself and wires
	// the Anthropic streaming iterator, which implements both __iter__ and
	// __aiter__.
	helpers.RequireCommandSucceeds(t, image, nil, "/app/.venv/bin/python", "-c", `
from litellm.llms.anthropic.chat.handler import ModelResponseIterator
from litellm.llms.opencode.chat.transformation import OpenCodeGoMessagesChatConfig

assert "get_model_response_iterator_for_request" in OpenCodeGoMessagesChatConfig.__dict__, (
    "override missing from OpenCodeGoMessagesChatConfig"
)
iterator = ModelResponseIterator(streaming_response=iter([]), sync_stream=True)
assert callable(iterator.__iter__) and callable(iterator.__aiter__)
print("override present")
`)

	// The OpenAI-wire config must forward messages without the per-message
	// reasoning fields that Go's /v1/chat/completions endpoint rejects.
	helpers.RequireCommandSucceeds(t, image, nil, "/app/.venv/bin/python", "-c", `
from litellm.llms.opencode.chat.transformation import OpenCodeGoChatConfig

body = OpenCodeGoChatConfig().transform_request(
    model="opencode_go/beam-free",
    messages=[
        {"role": "user", "content": "hi"},
        {"role": "assistant", "content": "ok", "reasoning": "chain-of-thought", "reasoning_content": "deep"},
        {"role": "user", "content": "again"},
    ],
    optional_params={"max_tokens": 100},
    litellm_params={"api_key": "sk-test"},
    headers={},
)
assert all("reasoning" not in m and "reasoning_content" not in m for m in body["messages"])
assert body["messages"][1]["content"] == "ok"
print("reasoning stripped")
`)

	// Every Go wire config must guarantee an x-opencode-session header: Go
	// answers 400 missing_session_id otherwise. Headerless callers get a
	// throwaway oneshot-<uuid> id; established sessions and direct client
	// headers must pass through untouched.
	helpers.RequireCommandSucceeds(t, image, nil, "/app/.venv/bin/python", "-c", `
from litellm.llms.opencode.chat.transformation import (
    OpenCodeGoChatConfig,
    OpenCodeGoMessagesChatConfig,
)

cfg = OpenCodeGoChatConfig()

def flat(h):
    return {str(k).lower(): v for k, v in h.items()}

def validate(cfg, headers, litellm_params):
    return flat(cfg.validate_environment(
        headers=headers,
        model="opencode_go/beam-free",
        messages=[],
        optional_params={},
        litellm_params=litellm_params,
        api_key="sk-test",
    ))

sess = validate(cfg, {}, {"api_key": "sk-test"}).get("x-opencode-session")
assert isinstance(sess, str) and sess.startswith("oneshot-"), sess

sess = validate(cfg, {}, {"api_key": "sk-test", "litellm_session_id": "sess-123"}).get("x-opencode-session")
assert sess == "sess-123", sess

sess = validate(cfg, {"x-opencode-session": "custom-sess"}, {"api_key": "sk-test"}).get("x-opencode-session")
assert sess == "custom-sess", sess

sess = validate(OpenCodeGoMessagesChatConfig(), {}, {"api_key": "sk-test"}).get("x-opencode-session")
assert isinstance(sess, str) and sess.startswith("oneshot-"), sess
print("session header guaranteed")
`)
}
