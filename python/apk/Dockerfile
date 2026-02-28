#syntax=docker/dockerfile:1.19

FROM container-base AS minimal-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

RUN --mount=type=bind,source=apk/minimal.yml,target=/minimal.yml apko build-minirootfs /minimal.yml /staging.tar.gz

WORKDIR /staging
RUN tar -xpf /staging.tar.gz -C /staging

FROM container-base AS dev-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

RUN --mount=type=bind,source=apk/dev.yml,target=/dev.yml apko build-minirootfs /dev.yml /staging.tar.gz

WORKDIR /staging
RUN tar -xpf /staging.tar.gz -C /staging

FROM container-base AS python-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

COPY --from=astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
RUN --mount=type=bind,source=install-python.sh,target=/install-python.sh <<-EOT
  mkdir -p /staging/usr/local/{bin,include,lib,share}
  sh /install-python.sh --python-version="${PYTHON_VERSION}" --path=/staging/usr/local --strip=true
EOT

FROM scratch AS minimal

COPY --from=minimal-builder /staging /
COPY --from=python-builder /staging/usr/local /usr/local

USER nonroot
ENTRYPOINT ["tini", "--", "python"]

FROM scratch AS dev

COPY --from=dev-builder /staging /
COPY --from=python-builder /staging/usr/local /usr/local
COPY --from=astral-sh/uv /uv /uvx /usr/local/bin/

ENTRYPOINT ["tini", "--"]
CMD ["python"]
