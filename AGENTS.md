# AGENTS.md: container-images

Guidance for AI coding agents (and humans) working in this repository.

This repo has no application code. Its purpose is to build and publish
container images to GHCR (`ghcr.io/cbc02009/<app>`); Go exists solely to test
those images via `testcontainers-go`.

## Layout

- `apps/<name>/Dockerfile` — the image. Keep base-image versions on their own
  line with a `# renovate: datasource=docker depName=<image>` comment directly
  above it so Renovate tracks it.
- `apps/<name>/docker-bake.hcl` — Bake targets. The `VERSION` variable is the
  single source of truth for the upstream version and must carry a
  `# renovate: datasource=... depName=...` comment directly above it.
  Targets: `image-local` (local, `type=docker`), `image` (inherits
  `docker-metadata-action`; used by CI), `image-all` (declares platforms).
- `apps/<name>/container_test.go` — optional, `package main`, one file per app.
- `tests/helpers.go` — shared test helpers (`package helpers`):
  `GetTestImage`, `RequireHTTPEndpoint`, `RequireExecSucceeds`,
  `RequireFileExists`, `RequireCommandSucceeds`. Check this file before
  hand-rolling container lifecycle code in an individual test; add a new
  capability here instead of duplicating it.

## Conventions

- `testify/require`, not `assert`: a failed image test should stop
  immediately rather than cascade into a second, confusing failure.
- Every helper takes `t *testing.T` first, calls `t.Helper()`, and registers
  cleanup via `testcontainers.CleanupContainer(t, c)`; never leak a container
  past the test.
- `TEST_IMAGE` overrides the image under test (`helpers.GetTestImage`), so a
  local build can point tests at a just-built image instead of a published
  tag.
- Tests are idempotent and side-effect-free: they only assert against the
  image under test (HTTP response, command exit code, file presence) and
  never depend on or mutate state from another test.
- Go must be `gofmt -s`-clean and `go vet`-clean (CI enforces both in the
  Pull Request workflow).
- Don't state a library's API from memory: verify against `pkg.go.dev` or
  this repo's own code before assuming a helper exists or behaves a certain
  way.

## Verification

After changing an app, run its test locally before calling the change done:

```sh
mise run local-build <app>
```

For workflow changes, run `zizmor .` (installed via mise) and re-read the
affected workflow. For Dockerfile changes, run `droast .`.

## Commits

- PR titles follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):
  `<type>[(scope)][!]: <description>`. Sign commits with `git commit -s`.
- Never `git push` or trigger the Release workflow unless asked; a push to
  `main` publishes images, which is a hard-to-reverse action.
