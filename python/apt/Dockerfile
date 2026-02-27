# syntax=docker/dockerfile:1.19

# ── Install Python via uv into a staging directory ──
FROM container-base AS python-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

COPY --from=astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
RUN --mount=type=bind,source=../install-python.sh,target=/install-python.sh <<-EOT
  mkdir -p /staging/usr/local/{bin,include,lib,share}
  sh /install-python.sh --python-version="${PYTHON_VERSION}" --path=/staging/usr/local --strip=false
EOT

# ── Common base: copy Python + create non-root user ──
FROM container-base AS python-base

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

COPY --from=python-builder /staging /

ARG USER_GID
ARG USER_UID
ARG USER_NAME
RUN <<-EOT
	groupadd --system --gid "${USER_GID}" -- "${USER_NAME}"
	useradd --create-home --system --uid "${USER_UID}" --gid "${USER_GID}" -- "${USER_NAME}"
EOT

ENTRYPOINT ["tini", "--"]
CMD ["python"]

# ── Minimal: matches python:slim runtime deps ──
FROM python-base AS minimal

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=tmpfs,target=/var/log <<-EOT
	apt-get update
	apt-get dist-upgrade -y
	apt-get install -y --no-install-recommends \
		ca-certificates \
		libexpat1 \
		media-types \
		netbase \
		tini \
		tzdata
	apt-get autoremove -y
	apt-get clean
EOT

ARG USER_NAME
USER ${USER_NAME}

# ── Dev: matches buildpack-deps + python build deps ──
FROM python-base AS dev

# buildpack-deps:bookworm-curl + scm + full + python-specific deps
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    --mount=type=tmpfs,target=/var/log <<-EOT
	apt-get update
	apt-get dist-upgrade -y
	apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
		gnupg \
		netbase \
		sq \
		wget \
		\
		git \
		mercurial \
		openssh-client \
		subversion \
		procps \
		\
		autoconf \
		automake \
		bzip2 \
		default-libmysqlclient-dev \
		dpkg-dev \
		file \
		g++ \
		gcc \
		imagemagick \
		libbz2-dev \
		libc6-dev \
		libcurl4-openssl-dev \
		libdb-dev \
		libevent-dev \
		libffi-dev \
		libgdbm-dev \
		libglib2.0-dev \
		libgmp-dev \
		libjpeg-dev \
		libkrb5-dev \
		liblzma-dev \
		libmagickcore-dev \
		libmagickwand-dev \
		libmaxminddb-dev \
		libncurses5-dev \
		libncursesw5-dev \
		libpng-dev \
		libpq-dev \
		libreadline-dev \
		libsqlite3-dev \
		libssl-dev \
		libtool \
		libwebp-dev \
		libxml2-dev \
		libxslt-dev \
		libyaml-dev \
		libzstd-dev \
		make \
		patch \
		unzip \
		xz-utils \
		zlib1g-dev \
		\
		libbluetooth-dev \
		tk-dev \
		uuid-dev \
		\
		tini
EOT

COPY --from=astral-sh/uv /uv /uvx /usr/local/bin/
