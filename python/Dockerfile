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
    - libexpat1
    - libffi
    - sqlite-libs
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

COPY --from=ghcr.io/astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
RUN <<-EOT
    mkdir -p /python
    uv python install ${PYTHON_VERSION} --install-dir=/python --no-bin
    cd python/*/
    for dir in bin include lib share; do mv $dir/* /rootfs/usr/local/$dir/; done

    # Remove unnecessary components
    rm -rf /rootfs/usr/local/bin/{pip*,idle*,pydoc*,python*-config} \
           /rootfs/usr/local/lib/{libtcl*.so*,tcl*,tk*,itcl*,thread*} \
           /rootfs/usr/local/lib/python*/{tkinter,idlelib,ensurepip,pydoc_data,turtle.py,turtledemo,__phello__} \
           /rootfs/usr/local/share/{man,terminfo}
EOT

FROM scratch

COPY --from=builder /rootfs /

ARG USER_NAME
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "python"]
