# syntax=docker/dockerfile:1.28
FROM golang AS build

SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

WORKDIR /app

ENV CGO_ENABLED=0

RUN --mount=type=bind,source=go.mod,target=/src/go.mod \
    --mount=type=bind,source=go.sum,target=/src/go.sum \
    --mount=type=bind,source=os-images,target=/src/os-images \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go -C /src build -ldflags="-s -w" -o /app/os-images ./os-images

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/os-images /os-images

STOPSIGNAL SIGTERM

USER 65532:65532
ENTRYPOINT ["/os-images"]
