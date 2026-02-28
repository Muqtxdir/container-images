# OS Base Images

Ubuntu and Debian base images, used as the foundation for other container images in this repository.

## Features

- **Upgraded packages**: `apt-get dist-upgrade` applied at build time
- **Process supervisor**: [tini](https://github.com/krallin/tini) as PID 1 via ENTRYPOINT
- **Non-root by default**: `nonroot` user pre-created (UID/GID: 65532)
- **Clean image**: Built via scratch copy to eliminate layer bloat
- **Modern compression**: Zstandard (zstd) compression for reduced image size
- **OCI compliant**: Full OCI image specification labels and annotations

## Tags

| Variant | Registry Tag |
| ------- | ------------ |
| Ubuntu  | `ghcr.io/muqtxdir/container-images/ubuntu:latest` |
| Debian  | `ghcr.io/muqtxdir/container-images/debian:latest` |

## Contributing

Please reach out at [github.com/muqtxdir/container-images](https://github.com/muqtxdir/container-images).
