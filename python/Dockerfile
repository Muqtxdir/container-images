# syntax=docker/dockerfile:1.19
FROM container-base AS builder

SHELL ["/bin/sh", "-euxo", "pipefail", "-c"]

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
	--mount=type=cache,target=/var/lib/apk/lists,sharing=locked \
	--mount=type=tmpfs,target=/var/log <<-EOT
    apk add apko
EOT

# Generate apko config dynamically for current architecture only
ARG PYTHON_VERSION
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
    - ca-certificates
    - python-${PYTHON_VERSION}
    - tini

accounts:
  groups:
    - groupname: ${USER_NAME}
      gid: ${USER_GID}
  users:
    - username: ${USER_NAME}
      uid: ${USER_UID}
      gid: ${USER_GID}
      home: /home/${USER_NAME}

paths:
  - path: /usr/local/bin
    type: directory
    permissions: 0o755
  - path: /usr/local/bin/python
    type: symlink
    source: /usr/bin/python3
  - path: /usr/local/bin/python3
    type: symlink
    source: /usr/bin/python3
  - path: /usr/local/bin/python${PYTHON_VERSION}
    type: symlink
    source: /usr/bin/python${PYTHON_VERSION}
EOF

# Build rootfs via apko
RUN <<-EOT
	apko build-minirootfs /base.apko.yml /rootfs.tar
	mkdir -p /rootfs
	tar -xf /rootfs.tar -C /rootfs
EOT


FROM scratch

COPY --link --from=builder /rootfs /

ARG USER_NAME
WORKDIR /home/${USER_NAME}
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "python"]
CMD ["--help"]
