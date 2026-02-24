# syntax=docker/dockerfile:1.19

# ── Fetch available Python versions from uv ──
FROM container-base AS uv-versions
COPY --from=astral-sh/uv /uv /usr/local/bin/uv
COPY <<'EOF' /fetch-versions.sh
#!/bin/sh
set -eu

# Fetch available CPython versions from uv and output a Docker Bake override JSON.
# Produces the latest patch per minor version in ascending order, plus a "latest" entry.
#
# Example output: ["3.8.20","3.9.25","3.10.19",...,"3.15.0a6","latest"]

# Step 1: Get cpython versions, skip freethreaded, extract version numbers
#   uv lists newest first — we rely on this ordering for dedup
uv python list --python-preference only-managed \
    | grep '^cpython-' \
    | grep -v '+freethreaded' \
    | cut -d- -f2 \
    > /tmp/all-versions.txt

# Step 2: Keep only the latest patch per minor (first occurrence = newest)
> /tmp/versions.txt
seen=""
while read -r version; do
    minor=$(echo "$version" | cut -d. -f1-2)
    echo "$seen" | grep -q "$minor" && continue
    seen="$seen $minor"
    echo "$version" >> /tmp/versions.txt
done < /tmp/all-versions.txt

# Step 3: Build bake override JSON with jq
jq -Rs '
    split("\n") | map(select(length > 0)) | reverse | . + ["latest"] |
    { variable: { PYTHON_VERSIONS_LIST: { default: . } } }
' /tmp/versions.txt
EOF
RUN sh /fetch-versions.sh > /versions-override.json

FROM scratch AS uv
COPY --from=uv-versions /versions-override.json /

# ── Build distroless Python image ──
FROM container-base AS builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

# Generate apko config dynamically for current architecture only
ARG USER_GID
ARG USER_UID
ARG USER_NAME
COPY <<EOF /base.apko.yml
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
    - libbz2-1
    - libcrypt1
    - libcrypto3
    - libexpat1
    - libffi
    - libgcc
    - libssl3
    - libstdc++
    - libuuid
    - libxcrypt
    - mpdecimal
    - ncurses
    - ncurses-terminfo-base
    - readline
    - sqlite-libs
    - xz
    - zlib
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

COPY --from=astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
RUN <<-EOT
    mkdir -p /python
    uv python install ${PYTHON_VERSION} --install-dir=/python --no-bin

    # Get the first matching Python installation directory
    python_dir=$(ls -d /python/cpython-* 2>/dev/null | head -n1)
    for dir in bin include lib share; do
        mv "$python_dir/$dir"/* /rootfs/usr/local/$dir/
    done

    # Remove unnecessary components
    rm -rf /rootfs/usr/local/bin/{pip*,idle*,pydoc*,python*-config} \
           /rootfs/usr/local/lib/{libtcl*.so*,tcl*,tk*,itcl*,thread*} \
           /rootfs/usr/local/lib/python*/{tkinter,idlelib,ensurepip,pydoc_data,turtle.py,turtledemo,__phello__} \
           /rootfs/usr/local/lib/python*/site-packages/{pip,pip-*,setuptools,setuptools-*,wheel,wheel-*} \
           /rootfs/usr/local/share/{man,terminfo}
EOT

FROM scratch

COPY --from=builder /rootfs /

ARG USER_NAME
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "python"]
