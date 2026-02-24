# syntax=docker/dockerfile:1.19

# ── Fetch available Python versions from uv ──
FROM container-base AS uv-versions
COPY --from=astral-sh/uv /uv /usr/local/bin/uv

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

COPY <<'EOF' /fetch-versions.sh
#!/bin/sh
set -eu


# Step 1: Get cpython versions, skip freethreaded, extract version numbers
uv python list --python-preference only-managed \
    | grep '^cpython-' \
    | grep -v '+freethreaded' \
    | cut -d- -f2 \
    | jq -Rs '
        # Split lines, drop empties
        split("\n") | map(select(length > 0)) |

        # Keep only the latest patch per minor (first = newest from uv)
        reduce .[] as $v ({seen: {}, out: []};
            ($v | split(".")[0:2] | join(".")) as $minor |
            if .seen[$minor] then .
            else .seen[$minor] = true | .out += [$v]
            end
        ) | .out |

        # Reverse to ascending order, append "latest" entry
        reverse | . + ["latest"] |

        # Wrap in bake override structure
        { variable: { PYTHON_VERSIONS_LIST: { default: . } } }
    '
EOF
RUN bash /fetch-versions.sh > /versions-override.json

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
