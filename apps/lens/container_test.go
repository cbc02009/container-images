package main

import (
	"testing"

	helpers "github.com/cbc02009/container-images/tests"
)

// Lens needs its ClickHouse sidecar (or an external ClickHouse) to bind its
// HTTP listener, so the full service cannot start in a bare container. The
// upstream Dockerfile's own sandbox smoke stage is likewise not built here.
// Instead assert that every runtime artifact shipped by the image exists:
// the Rust service and migration binaries, the exported Next.js UI, the
// agent-chat supervisor entrypoint, and the sandbox policy/runtime files.
func TestLens(t *testing.T) {
	image := helpers.GetTestImage("lens:v1.108.0-rc.1")

	helpers.RequireFileExists(t, image, "/usr/local/bin/litellm-lens")
	helpers.RequireFileExists(t, image, "/usr/local/bin/lens-migrate")
	helpers.RequireFileExists(t, image, "/app/ui/index.html")
	helpers.RequireFileExists(t, image, "/app/agent-chat/dist/supervisor.js")
	helpers.RequireFileExists(t, image, "/app/lens/python.seccomp")
	helpers.RequireFileExists(t, image, "/app/lens/python-runtime.json")
}
