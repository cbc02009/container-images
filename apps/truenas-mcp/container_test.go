package main

import (
	"testing"

	helpers "github.com/cbc02009/container-images/tests"
)

// The image wraps the official truenas/truenas-mcp binary (upstream ships no
// container image) on distroless static. The server only ever talks to a real
// TrueNAS over wss:// and exits at startup if it cannot reach one, so it
// cannot serve without a NAS; this test asserts the runtime contract instead:
// the binary is present, executable, and reports the pinned release version.
// The live cluster behaviour (version gate against TrueNAS 25.10+, tool
// inventory, elicitation/confirmation flow) is verified in k8s-home-ops after
// the toolhive server is deployed.
func TestTruenasMcp(t *testing.T) {
	image := helpers.GetTestImage("truenas-mcp:1.0.1")

	helpers.RequireFileExists(t, image, "/usr/local/bin/truenas-mcp")

	helpers.RequireCommandSucceeds(t, image, nil, "truenas-mcp", "--version")
}
