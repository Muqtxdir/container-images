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

EOF

# Build rootfs via apko and create symlinks
RUN <<-EOT
	apko build-minirootfs /base.apko.yml /rootfs.tar
	mkdir -p /rootfs
	tar -xpf /rootfs.tar -C /rootfs

	# Create symlinks in /usr/local/bin manually
	mkdir -p /rootfs/usr/local/bin
	ln -sf /usr/bin/python3 /rootfs/usr/local/bin/python
	ln -sf /usr/bin/python3 /rootfs/usr/local/bin/python3
	ln -sf /usr/bin/python${PYTHON_VERSION} /rootfs/usr/local/bin/python${PYTHON_VERSION}
EOT


FROM scratch

COPY --link --from=builder /rootfs /

ARG USER_NAME
USER ${USER_NAME}

ENTRYPOINT ["tini", "--", "python"]