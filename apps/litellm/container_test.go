package main

import (
	"testing"

	helpers "github.com/cbc02009/container-images/tests"
)

// The image is an overlay on the upstream LiteLLM image carrying the
// opencode_go /v1/messages streaming fix: OpenCodeGoMessagesChatConfig must
// define get_model_response_iterator_for_request itself (the inherited
// BaseConfig stub returns None, which kills OpenAI-wire streaming to
// /v1/messages Go models with 'NoneType' object is not an iterator).
func TestLitellm(t *testing.T) {
	image := helpers.GetTestImage("litellm:v1.109.0-rc.1")

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
}
