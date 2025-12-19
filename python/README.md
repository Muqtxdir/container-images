# Python Images

Minimal, Python container images built from scratch using [Chainguard Wolfi](https://github.com/wolfi-dev) packages.

## Features

- **Minimal footprint**: Distroless design with only essential components
- **Security-focused**: Runs as non-root user (UID/GID: 65532)
- **Modern compression**: Uses Zstandard (zstd) compression for reduced image size
- **OCI compliant**: Full OCI image specification labels and annotations

## Available Versions

- Python 3.11
- Python 3.12
- Python 3.13
- Python 3.14

## Contributing

Contributions are welcome! Please open an issue or pull request.

## References

- [Chainguard Wolfi](https://github.com/wolfi-dev)
- [OCI Image Specification](https://github.com/opencontainers/image-spec)
- [Docker Buildx Bake](https://docs.docker.com/build/bake/)
