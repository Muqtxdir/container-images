# syntax=docker/dockerfile:1.28
FROM golang AS build

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

WORKDIR /app

ENV CGO_ENABLED=0

RUN --mount=type=bind,source=go.mod,target=/src/go.mod \
    --mount=type=bind,source=go.sum,target=/src/go.sum \
    --mount=type=bind,source=internal,target=/src/internal \
    --mount=type=bind,source=os-images,target=/src/os-images \
    --mount=type=bind,source=python-images,target=/src/python-images \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go -C /src build -ldflags="-s -w" -o /app/ ./os-images ./python-images

FROM scratch AS base
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

STOPSIGNAL SIGTERM

FROM base AS os-images
COPY --from=build /app/os-images /os-images
ENTRYPOINT ["/os-images"]
