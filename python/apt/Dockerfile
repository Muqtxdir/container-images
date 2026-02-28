#syntax=docker/dockerfile:1.19

FROM container-base AS python-builder

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

COPY --from=astral-sh/uv /uv /bin/uv

ARG PYTHON_VERSION
ARG STRIP_EXTRAS
RUN --mount=type=bind,source=install-python.sh,target=/install-python.sh <<-EOT
  	mkdir -p /staging/usr/local/{bin,include,lib,share}
  	sh /install-python.sh --python-version="${PYTHON_VERSION}" --path=/staging/usr/local --strip=${STRIP_EXTRAS}
EOT

FROM container-base AS minimal

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
		tzdata
	apt-get autoremove -y
	apt-get clean
EOT

COPY --from=python-builder /staging /

USER nonroot
CMD ["python"]

FROM container-base AS dev

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
		uuid-dev
EOT

COPY --from=astral-sh/uv /uv /uvx /usr/local/bin/
COPY --from=python-builder /staging /

CMD ["python"]
