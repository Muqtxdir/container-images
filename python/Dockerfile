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
    mkdir -p /rootfs
    tar -xpf /rootfs.tar -C /rootfs
    mkdir -p /rootfs/usr/local/bin /rootfs/usr/local/lib /rootfs/usr/local/include
EOT

FROM python-base AS python

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

RUN <<-EOT
    pip uninstall -y pip setuptools wheel
    rm -rf /usr/local/lib/python*/site-packages/pip* \
           /usr/local/lib/python*/site-packages/setuptools* \
           /usr/local/lib/python*/site-packages/wheel* \
           /usr/local/bin/pip*
EOT

FROM scratch

COPY --from=builder /rootfs /
COPY --from=python /usr/local/ /usr/local/

ARG USER_NAME
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "python"]