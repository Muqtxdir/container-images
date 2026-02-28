variable "GITHUB_SHA" {
    default = "main"
}

variable "GITHUB_REPOSITORY" {
    default = "muqtxdir/container-images"
}

variable "GITHUB_ACTOR" {
    default = "Muqtxdir"
}

variable "GITHUB_RUN_ID" {
    default = "local"
}

# ---------------------------------------------------------------------------
# OCI Labels+Annotations functions
# ---------------------------------------------------------------------------

function "oci_labels" {
    params = [distro]
    result = {
      "org.opencontainers.image.created"       =  "${timestamp()}"
      "org.opencontainers.image.authors"        = "${GITHUB_ACTOR}"
      "org.opencontainers.image.url"            = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.documentation"  = "https://github.com/${GITHUB_REPOSITORY}/blob/main/os/README.md"
      "org.opencontainers.image.source"         = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.version"        = "latest"
      "org.opencontainers.image.revision"       = "${GITHUB_SHA}"
      "org.opencontainers.image.vendor"         = "${GITHUB_ACTOR}"
      "org.opencontainers.image.licenses"       = "GPL-3.0-or-later"
      "org.opencontainers.image.title"          = "${distro}"
      "org.opencontainers.image.description"    = "${distro} image with tini and nonroot user."
      "org.opencontainers.image.base.name"      = "${distro}:latest"
    }
}

function "oci_annotations" {
    params = [distro]
    result = flatten([
        for scope in ["index", "manifest"] : [
            "${scope}:org.opencontainers.image.created=${timestamp()}",
            "${scope}:org.opencontainers.image.authors=${GITHUB_ACTOR}",
            "${scope}:org.opencontainers.image.url=https://github.com/${GITHUB_REPOSITORY}",
            "${scope}:org.opencontainers.image.documentation=https://github.com/${GITHUB_REPOSITORY}/blob/main/os/README.md",
            "${scope}:org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}",
            "${scope}:org.opencontainers.image.version=latest",
            "${scope}:org.opencontainers.image.revision=${GITHUB_SHA}",
            "${scope}:org.opencontainers.image.vendor=${GITHUB_ACTOR}",
            "${scope}:org.opencontainers.image.licenses=GPL-3.0-or-later",
            "${scope}:org.opencontainers.image.title=${distro}",
            "${scope}:org.opencontainers.image.description=${distro} image with tini and nonroot user.",
            "${scope}:org.opencontainers.image.base.name=${distro}:latest"
        ]
    ])
}

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

target "default" {
    platforms = ["linux/arm64", "linux/amd64"]
    output    = ["type=image,compression=zstd,compression-level=19,oci-mediatypes=true,force-compression=true"]
    pull      = true
    attest = [
        "type=provenance,mode=max,builder-id=https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}",
        "type=sbom,generator=docker/buildkit-syft-scanner"
    ]
    name = "${distro}"
    matrix = {
        distro = ["ubuntu", "debian"]
    }
    contexts = {
      "container-base" = "docker-image://${distro == "ubuntu" ? "public.ecr.aws/docker/library/ubuntu:latest" : "public.ecr.aws/docker/library/debian:latest"}"
    }
    dockerfile = "Dockerfile"
    target = "os"
    args = {
        USER_GID  = "65532"
        USER_UID  = "65532"
        USER_NAME = "nonroot"
    }
    tags        = ["ghcr.io/${lower(GITHUB_REPOSITORY)}/${distro}:latest"]
    labels      = oci_labels(distro)
    annotations = oci_annotations(distro)
    cache-from  = ["type=registry,ref=ghcr.io/${lower(GITHUB_REPOSITORY)}/buildcache:${distro}"]
    cache-to    = ["type=registry,ref=ghcr.io/${lower(GITHUB_REPOSITORY)}/buildcache:${distro},mode=max,compression=zstd,compression-level=12,ignore-error=true"]
}
