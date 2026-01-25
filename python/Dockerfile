# syntax=docker/dockerfile:1.19
FROM container-base AS builder

SHELL ["/bin/sh", "-euxo", "pipefail", "-c"]

# Generate apko config dynamically for current architecture only
ARG USER_GID
ARG USER_UID
ARG USER_NAME
COPY <<-EOF /base.apko.yml
contents:
  keyring:
    - https://packages.wolfi.dev/os/wolfi-signing.rsa.pub
  repositories:
    - https://packages.wolfi.dev/os
  packages:
    - ca-certificates-bundle
    - glibc
    - glibc-locale-posix
    - ld-linux
    - libgcc
    - libstdc++
    - libcrypto3
    - libssl3
    - zlib
    - xz
    - libbz2-1
    - libexpat1
    - libffi
    - readline
    - ncurses
    - sqlite-libs
    - mpdecimal
    - tzdata
    - tini

accounts:
  groups:
    - groupname: ${USER_NAME}
      gid: ${USER_GID}
  users:
    - username: ${USER_NAME}
      uid: ${USER_UID}
      gid: ${USER_GID}

EOF

# Build rootfs via apko and create symlinks
RUN <<-EOT
    apko build-minirootfs /base.apko.yml /rootfs.tar
    mkdir -p /rootfs/usr/local/{bin,include,lib,share}
    tar -xpf /rootfs.tar -C /rootfs
EOT

COPY --from=ghcr.io/astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
WORKDIR /{PYTHON_VERSION}
RUN <<-EOT
    uv python install {PYTHON_VERSION} --install-dir=/{PYTHON_VERSION} --no-bin
    for dir in bin include lib share; do mv $dir/* /rootfs/usr/local/$dir/; done
EOT

FROM scratch

COPY --from=builder /rootfs /

ARG USER_NAME
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "python"]
