# Node.js Images

Minimal, Node.js container images built from scratch using [official Node.js binaries](https://nodejs.org) and [Chainguard Wolfi](https://github.com/wolfi-dev) packages.

## Features

- **Minimal footprint**: Distroless design with only essential components
- **Security-focused**: Runs as non-root user (UID/GID: 65532)
- **Modern compression**: Uses Zstandard (zstd) compression for reduced image size
- **OCI compliant**: Full OCI image specification labels and annotations
- **Runtime-only**: No package managers (npm, npx, corepack) - intended for multi-stage builds

## Available Versions

| Version | Status | Docker Hub Tag | GitHub Tag |
|---------|--------|----------------|------------|
| 25 | Current | `muqtxdir/node:25` | `ghcr.io/muqtxdir/container-images/node:25` |
| 24 | LTS (Krypton) | `muqtxdir/node:24` | `ghcr.io/muqtxdir/container-images/node:24` |
| 22 | LTS (Jod) | `muqtxdir/node:22` | `ghcr.io/muqtxdir/container-images/node:22` |
| 20 | LTS (Iron) | `muqtxdir/node:20` | `ghcr.io/muqtxdir/container-images/node:20` |

**Note:** Node.js 25 is also available as `latest` tag. Node.js 24 is also available as `lts` tag on both registries.

## Usage

### Pull Images

- Pull the latest Node.js Image via `docker`
    ```
    docker image pull muqtxdir/node
    ```

- Pull the latest Node.js Image via `podman`
    ```
    podman image pull muqtxdir/node
    ```

## Contributing

Please reach out at [github.com/muqtxdir/container-images](https://github.com/muqtxdir/container-images).

## References

- [Node.js Official Downloads](https://nodejs.org/en/download)
- [Node.js Release Schedule](https://github.com/nodejs/release#release-schedule)
- [Chainguard Wolfi](https://github.com/wolfi-dev)
- [OCI Image Specification](https://github.com/opencontainers/image-spec)
- [Docker Buildx Bake](https://docs.docker.com/build/bake/)
