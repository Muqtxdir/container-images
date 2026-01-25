# Python Images

Minimal, Python container images built from scratch using [uv's standalone Python](https://github.com/astral-sh/python-build-standalone) and [Chainguard Wolfi](https://github.com/wolfi-dev) packages.

## Features

- **Minimal footprint**: Distroless design with only essential components
- **Security-focused**: Runs as non-root user (UID/GID: 65532)
- **Modern compression**: Uses Zstandard (zstd) compression for reduced image size
- **OCI compliant**: Full OCI image specification labels and annotations

## Available Versions

| Version | Status | EOL Date | Docker Hub Tag | GitHub Tag |
|---------|--------|----------|----------------|------------|
| 3.14 | Supported | October 2030 | `muqtxdir/python:3.14` | `ghcr.io/muqtxdir/container-images/python:3.14` |
| 3.13 | Supported | October 2029 | `muqtxdir/python:3.13` | `ghcr.io/muqtxdir/container-images/python:3.13` |
| 3.12 | Supported | October 2028 | `muqtxdir/python:3.12` | `ghcr.io/muqtxdir/container-images/python:3.12` |
| 3.11 | Supported | October 2027 | `muqtxdir/python:3.11` | `ghcr.io/muqtxdir/container-images/python:3.11` |
| 3.10 | Supported | October 2026 | `muqtxdir/python:3.10` | `ghcr.io/muqtxdir/container-images/python:3.10` |
| 3.9 | End of Life | October 2025 | `muqtxdir/python:3.9` | `ghcr.io/muqtxdir/container-images/python:3.9` |
| 3.8 | End of Life | October 2024 | `muqtxdir/python:3.8` | `ghcr.io/muqtxdir/container-images/python:3.8` |

**Note:** Python 3.14 is also available as `latest` and `3` tags on both registries.

## Usage

- Pull the latest Python Image via `docker`
    ```
    docker image pull muqtxdir/python
    ```

- Pull the latest Python Image via `podman`
    ```
    podman image pull muqtxdir/python
    ```

## Contributing

Please reach out at [github.com/muqtxdir/container-images](https://github.com/muqtxdir/container-images).

## References

- [python-build-standalone](https://github.com/indygreg/python-build-standalone)
- [uv - Python package manager](https://github.com/astral-sh/uv)
- [Chainguard Wolfi](https://github.com/wolfi-dev)
- [OCI Image Specification](https://github.com/opencontainers/image-spec)
- [Docker Buildx Bake](https://docs.docker.com/build/bake/)
