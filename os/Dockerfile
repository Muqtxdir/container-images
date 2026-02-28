FROM container-base AS os-builder

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=tmpfs,target=/var/log <<-EOT
	apt-get update
	apt-get dist-upgrade -y
    apt-get install -y --no-install-recommends tini
    apt-get autoremove
    apt-get clean
EOT

ARG USER_GID
ARG USER_UID
ARG USER_NAME
RUN <<-EOT
	groupadd --system --gid "${USER_GID}" -- "${USER_NAME}"
	useradd --create-home --system --uid "${USER_UID}" --gid "${USER_GID}" -- "${USER_NAME}"
EOT

FROM scratch AS os

COPY --from=os-builder / /

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

ENTRYPOINT [ "tini", "--" ]
CMD ["/bin/bash"]
