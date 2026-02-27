# Python Images

Python container images built using [uv's standalone Python](https://github.com/astral-sh/python-build-standalone), available across multiple base image variants.

## Features

- **Multiple bases**: Wolfi, Debian, and Ubuntu variants
- **Default & Dev variants**: Python runtime images and full development images with build tools
- **Security-focused**: Runs as non-root user (UID/GID: 65532)
- **Modern compression**: Uses Zstandard (zstd) compression for reduced image size
- **OCI compliant**: Full OCI image specification labels and annotations

## Tags
| Variant | DockerHub Tags | Github Registry Tags |
| ------- | ---------------| ---------------------|
| Wolfi/Chainguard | `muqtxdir/python:{version}` | `ghcr.io/muqtxdir/container-images/python:{version}` |
| Wolfi/Chainguard Dev | `muqtxdir/python:{version}-dev` | `ghcr.io/muqtxdir/container-images/python:{version}-dev` |
| Debian | `muqtxdir/python:{version}-debian` | `ghcr.io/muqtxdir/container-images/python:{version}-debian` |
| Debian Dev | `muqtxdir/python:{version}-debian-dev` | `ghcr.io/muqtxdir/container-images/python:{version}-debian-dev` |
| Ubuntu | `muqtxdir/python:{version}-ubuntu` | `ghcr.io/muqtxdir/container-images/python:{version}-ubuntu` |
| Ubuntu Dev | `muqtxdir/python:{version}-ubuntu-dev` | `ghcr.io/muqtxdir/container-images/python:{version}-ubuntu-dev` |

## Differences
| Variant | Default Tags | Dev Tags |
| --------| -------------| ---------|
| Wolfi/Chainguard | Replicates Chainguard's `python` tags | Replicates Chainguard's `python` dev tags |
| Debian | Replicates `python` slim tags | Replicates `python` deafult tags and includes uv |
| Ubuntu | Same as Debian but based on `ubuntu:latest` | Same as Debian Dev but based on `ubuntu:latest`|

## Contributing

Please reach out at [github.com/muqtxdir/container-images](https://github.com/muqtxdir/container-images).

## References

- [python-build-standalone](https://github.com/indygreg/python-build-standalone)
- [uv - Python package manager](https://github.com/astral-sh/uv)
- [Chainguard Wolfi](https://github.com/wolfi-dev)
- [OCI Image Specification](https://github.com/opencontainers/image-spec)
- [Docker Buildx Bake](https://docs.docker.com/build/bake/)
