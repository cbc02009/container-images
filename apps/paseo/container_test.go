package main

import (
	"testing"

	helpers "github.com/cbc02009/container-images/tests"
)

// The image starts the paseo daemon headless on boot, so a healthy
// container serves the API health endpoint and carries the agent CLI plus
// the home-ops toolchain.
func TestPaseo(t *testing.T) {
	image := helpers.GetTestImage("paseo:0.15.4")

	c := helpers.RequireHTTPEndpoint(t, image, helpers.HTTPTestConfig{
		Port: "6767",
		Path: "/api/health",
	}, nil)

	// The coding agent the daemon spawns.
	helpers.RequireExecSucceeds(t, c, "opencode", "--version")

	// The GitOps toolchain installed by mise.
	helpers.RequireExecSucceeds(t, c, "kubectl", "version", "--client")
	helpers.RequireExecSucceeds(t, c, "flux", "version", "--client")
	helpers.RequireExecSucceeds(t, c, "kustomize", "version")
	helpers.RequireExecSucceeds(t, c, "helm", "version")
	helpers.RequireExecSucceeds(t, c, "sops", "--version")
	helpers.RequireExecSucceeds(t, c, "age", "--version")
	helpers.RequireExecSucceeds(t, c, "talosctl", "version")
	helpers.RequireExecSucceeds(t, c, "topf", "--version")
	helpers.RequireExecSucceeds(t, c, "just", "--version")
	helpers.RequireExecSucceeds(t, c, "gh", "--version")
}
