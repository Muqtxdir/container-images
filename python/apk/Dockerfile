# syntax=docker/dockerfile:1.19

FROM container-base AS minimal-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

RUN --mount=type=bind,source=apk/minimal.yml,target=/minimal.yml <<-EOT
    apko build-minirootfs /minimal.yml /rootfs.tar.gz
    mkdir -p /rootfs/usr/local/{bin,include,lib,share}
    tar -xpf /rootfs.tar.gz -C /rootfs
EOT

COPY --from=astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
RUN --mount=type=bind,source=install-python.sh,target=/install-python.sh <<-EOT
  sh /install-python.sh --python-version="${PYTHON_VERSION}" --path=/rootfs/usr/local --strip=true
EOT

FROM container-base AS dev-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

RUN --mount=type=bind,source=apk/dev.yml,target=/dev.yml <<-EOT
    apko build-minirootfs /dev.yml /rootfs.tar.gz
    mkdir -p /rootfs/usr/local/{bin,include,lib,share}
    tar -xpf /rootfs.tar.gz -C /rootfs
EOT

COPY --from=astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
RUN --mount=type=bind,source=install-python.sh,target=/install-python.sh <<-EOT
  sh /install-python.sh --python-version="${PYTHON_VERSION}" --path=/rootfs/usr/local --strip=false
EOT

FROM scratch AS minimal

COPY --from=minimal-builder /rootfs /

USER nonroot
ENTRYPOINT ["tini", "--", "python"]

FROM scratch AS dev

COPY --from=dev-builder /rootfs /
COPY --from=astral-sh/uv /uv /uvx /usr/local/bin/

ENTRYPOINT ["tini", "--"]
CMD ["python"]
