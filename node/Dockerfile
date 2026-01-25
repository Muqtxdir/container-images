# syntax=docker/dockerfile:1.19
FROM container-base AS builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

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
    - ncurses-terminfo-base
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
    apko build-minirootfs /base.apko.yml /rootfs.tar.gz
    mkdir -p /rootfs/usr/local/{bin,include,lib,share}
    tar -xpf /rootfs.tar.gz -C /rootfs
EOT

ARG NODE_VERSION
RUN <<-EOT
    # Determine architecture for Node.js download
    ARCH=$(uname -m)
    if [ "$ARCH" = "x86_64" ]; then NODE_ARCH="x64"; fi
    if [ "$ARCH" = "aarch64" ]; then NODE_ARCH="arm64"; fi

    # Parse SHASUMS256.txt to find the latest version for this architecture
    LATEST_VERSION=$(curl -fsSL "https://nodejs.org/dist/latest-v${NODE_VERSION}.x/SHASUMS256.txt" \
        | grep -oE "node-v${NODE_VERSION}\.[0-9]+\.[0-9]+-linux-${NODE_ARCH}\.tar\.xz" \
        | sed "s/node-v\(.*\)-linux-${NODE_ARCH}\.tar\.xz/\1/" \
        | sort -V \
        | tail -1)

    # Download and extract official Node.js binary
    curl -fsSL "https://nodejs.org/dist/latest-v${NODE_VERSION}.x/node-v${LATEST_VERSION}-linux-${NODE_ARCH}.tar.xz" \
        | tar -xJ --strip-components=1 -C /rootfs/usr/local/

    # Remove package managers and unnecessary components
    rm -rf /rootfs/usr/local/lib/node_modules/{npm,corepack} \
           /rootfs/usr/local/bin/{npm,npx,corepack} \
           /rootfs/usr/local/share/{man,doc,systemtap} \
           /rootfs/usr/local/include/node/openssl
EOT

FROM scratch

COPY --from=builder /rootfs /

ARG USER_NAME
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "node"]
