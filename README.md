# container-images

A monorepo of Dockerfiles that builds and publishes container images to the
[GitHub Container Registry](https://ghcr.io) (`ghcr.io/cbc02009/<app>`), with
automatic upstream version updates via
[Renovate](https://github.com/renovatebot/renovate).

## How it works

Each image lives in `apps/<name>/` and consists of:

| File               | Purpose                                                                       |
| ------------------ | ----------------------------------------------------------------------------- |
| `Dockerfile`       | The image definition. Pin versions via `# renovate:` comments.                |
| `docker-bake.hcl`  | Bake targets: `image-local` (local dev), `image` (CI), `image-all` (platforms). The `VERSION` variable is the single source of truth for the upstream version, with a `# renovate:` annotation. |
| `container_test.go`| Optional Go test (testcontainers) that runs against a freshly built image in PRs. |

### Releases (`main` branch)

When a push to `main` changes anything under `apps/**`, the
`Release` workflow builds every changed app:

1. **plan** — reads `VERSION`/`SOURCE`/platforms from `docker-bake.hcl`,
   normalizes the version, and computes tags (`x.y.z`, `x.y`, `x`, raw
   version, `rolling`).
2. **build** — one job per platform. Builds with Buildx using registry layer
   caching (`ghcr.io/cbc02009/build_cache:<app>-<arch>`), pushes per-platform
   images by digest, and (for PRs only) runs the app's Go tests against a
   `sandbox`-tagged image.
3. **merge** — stitches per-platform images into one multi-arch manifest with
   all computed tags.
4. **attest** — uploads an SBOM and pushes a SLSA provenance attestation to
   the registry, then verifies it.

### Pull requests

The `Pull Request` workflow builds any changed apps **without publishing**
(tagged `sandbox` only, never pushed to GHCR) and runs its Go tests, plus
`go vet`/gofmt, [droast](https://github.com/immanuwell/dockerfile-roast)
Dockerfile linting, and [zizmor](https://zizmor.sh) workflow linting.

### Version updates (Renovate)

Renovate runs hourly using a GitHub App. It updates:

- the `VERSION` variable in `docker-bake.hcl` (upstream releases),
- `# renovate:`-annotated base images in `Dockerfile`s.

Updates to `docker-bake.hcl` auto-merge; the resulting `main` push triggers a
new release build automatically.

## Adding a new app

1. Create `apps/<name>/` with a `Dockerfile` and `docker-bake.hcl`
   (copy `apps/swarmui` as a template; the `VERSION` variable needs a
   `# renovate:` annotation pointing at the upstream release source).
2. Optionally add `container_test.go` (use the helpers in `tests/helpers.go`;
   see `AGENTS.md` for conventions).
3. Commit and push to `main` — the image is built and published as
   `ghcr.io/cbc02009/<name>`.

## Local development

Tooling is managed by [mise](https://mise.jdx.dev/) (`.mise/config.toml`).
After cloning:

```sh
mise install          # installs go, docker buildx deps, droast, zizmor, lefthook, ...
```

- `mise run local-build <app>` — builds the app's image locally and runs its
  Go tests against it.
- `mise run remote-build <app> [release]` — triggers the Release workflow
  from the command line.

`lefthook` (installed by mise's postinstall hook) lints staged Dockerfiles
(droast), Go files (gofmt), and shell scripts (shellcheck) on commit.

## Repository secrets

| Secret               | Used by      | Notes                                                             |
| -------------------- | ------------ | ----------------------------------------------------------------- |
| `BOT_APP_ID`          | Renovate        | GitHub App ID (the App needs Contents + Pull Requests: read/write and must be installed on this account). |
| `BOT_APP_PRIVATE_KEY` | Renovate        | GitHub App private key (PEM; PKCS#1 or PKCS#8, normalized automatically). |

No other secrets are required: GHCR pushes use the built-in `GITHUB_TOKEN`.
