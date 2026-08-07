# Releasing DUG

This document describes how to cut a release of DUG using [GoReleaser](https://goreleaser.com/).

## Prerequisites

- Go 1.25+ (see `go.mod`)
- [GoReleaser v2](https://goreleaser.com/install/) installed locally
- [Docker](https://docs.docker.com/get-docker/) with Buildx (for container images)
- A clean git working tree on the release commit
- An annotated semver tag (for example `v1.0.0`)

## Registry placeholders

Before publishing Docker images, replace the registry users in `.goreleaser.yaml` **or** set these environment variables (CI sets them automatically):

| Variable | Purpose | Example |
|----------|---------|---------|
| `GITHUB_USER` | GitHub username or org; used for `ghcr.io/<user>/dug` | `my-org` |
| `DOCKERHUB_USER` | Docker Hub username; used for `docker.io/<user>/dug` | `my-dockerhub-user` |

For local snapshot builds, export placeholder values if you have not customized the config:

```bash
export GITHUB_USER=your-github-user
export DOCKERHUB_USER=your-dockerhub-user
```

## Required GitHub secrets

Configure these in **Settings → Secrets and variables → Actions**:

| Secret | Required for | Description |
|--------|--------------|-------------|
| `GITHUB_TOKEN` | GitHub Releases, GHCR | Provided automatically by GitHub Actions; needs `contents: write` and `packages: write` |
| `DOCKERHUB_USERNAME` | Docker Hub | Your Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub | Docker Hub access token ([create one](https://hub.docker.com/settings/security)) |

The release workflow maps:

- `GITHUB_USER` → `github.repository_owner`
- `DOCKERHUB_USER` → `secrets.DOCKERHUB_USERNAME`

## What a release produces

GoReleaser builds and publishes:

- Cross-platform binaries for `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, and `windows/amd64`
- `.tar.gz` archives (Linux/macOS) and `.zip` archives (Windows) including `LICENSE` and `README.md`
- `SHA256SUMS.txt` checksum file
- GitHub Release with auto-generated release notes from commits
- Multi-arch Docker images (`linux/amd64`, `linux/arm64`) tagged as `latest` and `{{ .Version }}`:
  - `ghcr.io/<GITHUB_USER>/dug`
  - `docker.io/<DOCKERHUB_USER>/dug`

Version metadata (`Version`, `Commit`, `BuildDate`, `Go`) is injected at link time. Run `dug version` to verify.

## Test locally

Validate the GoReleaser configuration:

```bash
export GOVERSION="$(go env GOVERSION)"
goreleaser check
```

Build a snapshot release (no publish, artifacts land in `dist/`):

```bash
export GOVERSION="$(go env GOVERSION)"
export GITHUB_USER=your-github-user      # replace with your GitHub user/org
export DOCKERHUB_USER=your-dockerhub-user # replace with your Docker Hub user
goreleaser release --snapshot --clean
```

Inspect outputs:

```bash
ls -la dist/
./dist/dug_linux_amd64_v1/dug version
cat dist/SHA256SUMS.txt
```

Snapshot Docker images are built locally (not pushed) when Docker Buildx is available. Multi-arch image builds on `linux/amd64` hosts require QEMU binfmt:

```bash
docker run --privileged --rm tonistiigi/binfmt --install all
```

To skip Docker during local testing:

```bash
goreleaser release --snapshot --clean --skip=publish
```

Or use the Makefile helpers:

```bash
make release-check
make release-snapshot
```

## Create a release

### Option A — GitHub Actions (recommended)

1. Ensure tests pass on `main`.
2. Configure the GitHub secrets listed above.
3. Create and push a semver tag:

   ```bash
   git tag -a v1.0.0 -m "v1.0.0"
   git push origin v1.0.0
   ```

4. The [Release workflow](.github/workflows/release.yml) runs tests, lint, and GoReleaser.
5. Verify the GitHub Release, checksums, and container registries.

### Option B — Manual release

1. Export credentials:

   ```bash
   export GITHUB_TOKEN="<personal-access-token-with-repo-scope>"
   export GITHUB_USER="<your-github-user>"
   export DOCKERHUB_USER="<your-dockerhub-user>"
   export GOVERSION="$(go env GOVERSION)"
   ```

2. Log in to registries:

   ```bash
   echo "$GITHUB_TOKEN" | docker login ghcr.io -u "$GITHUB_USER" --password-stdin
   echo "$DOCKERHUB_TOKEN" | docker login -u "$DOCKERHUB_USER" --password-stdin
   ```

3. Tag and release:

   ```bash
   git tag -a v1.0.0 -m "v1.0.0"
   goreleaser release --clean
   git push origin v1.0.0
   ```

## Publish Docker images only

Docker images are built and pushed during the publish phase of `goreleaser release`. A normal tagged release handles this automatically.

To test Docker image builds without publishing:

```bash
export GOVERSION="$(go env GOVERSION)"
goreleaser release --snapshot --clean
```

Pull published images after a release:

```bash
docker pull ghcr.io/<GITHUB_USER>/dug:v1.0.0
docker pull docker.io/<DOCKERHUB_USER>/dug:v1.0.0
```

## Troubleshooting

- **`goreleaser check` fails** — ensure `.goreleaser.yaml` is valid YAML and GoReleaser v2 is installed.
- **Docker build fails locally** — install Buildx and ensure the Docker daemon is running; use `--skip=publish` to test binaries only.
- **GHCR push denied** — confirm workflow `permissions` include `packages: write`.
- **Docker Hub push denied** — verify `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets.
