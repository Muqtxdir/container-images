variable "GITHUB_SHA" {
    default = "main"
}

variable "GITHUB_REPOSITORY" {
    default = "muqtxdir/container-images"
}

variable "DOCKER_REPOSITORY" {
    default = "muqtxdir"
}

variable "GITHUB_ACTOR" {
    default = "Muqtxdir"
}

variable "GITHUB_RUN_ID" {
    default = "local"
}

variable "PYTHON_VERSIONS_LIST" {
    default = ["3.13", "3.14", "latest"]
}

# ---------------------------------------------------------------------------
# Tag functions
# ---------------------------------------------------------------------------

function "python_image_tags" {
    params = [version, suffix, github_repo, docker_repo]
    result = version == "latest" ? [
        "ghcr.io/${lower(github_repo)}/python:latest${suffix}",
        "ghcr.io/${lower(github_repo)}/python:3${suffix}",
        "docker.io/${lower(docker_repo)}/python:latest${suffix}",
        "docker.io/${lower(docker_repo)}/python:3${suffix}"
    ] : length(regexall("[a-zA-Z]", version)) > 0 ? [
        "ghcr.io/${lower(github_repo)}/python:${version}${suffix}",
        "docker.io/${lower(docker_repo)}/python:${version}${suffix}"
    ] : [
        "ghcr.io/${lower(github_repo)}/python:${version}${suffix}",
        "docker.io/${lower(docker_repo)}/python:${version}${suffix}",
        "ghcr.io/${lower(github_repo)}/python:${join(".", slice(split(".", version), 0, 2))}${suffix}",
        "docker.io/${lower(docker_repo)}/python:${join(".", slice(split(".", version), 0, 2))}${suffix}"
    ]
}

function "oci_labels" {
    params = [version, suffix]
    result = {
      "org.opencontainers.image.created"       = "${timestamp()}"
      "org.opencontainers.image.authors"        = "${GITHUB_ACTOR}"
      "org.opencontainers.image.url"            = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.documentation"  = "https://github.com/${GITHUB_REPOSITORY}/blob/main/python/README.md"
      "org.opencontainers.image.source"         = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.version"        = "${version}"
      "org.opencontainers.image.revision"       = "${GITHUB_SHA}"
      "org.opencontainers.image.vendor"         = "${GITHUB_ACTOR}"
      "org.opencontainers.image.licenses"       = "GPL-3.0-or-later"
      "org.opencontainers.image.title"          = "Python ${version}${suffix}"
      "org.opencontainers.image.description"    = "Python ${version}${suffix} image."
      "org.opencontainers.image.base.name"      = suffix == "-dev" ? "container-base" : "scratch"
    }
}

function "oci_annotations" {
    params = [version, suffix]
    result = flatten([
        for scope in ["index", "manifest"] : [
            "${scope}:org.opencontainers.image.created=${timestamp()}",
            "${scope}:org.opencontainers.image.authors=${GITHUB_ACTOR}",
            "${scope}:org.opencontainers.image.url=https://github.com/${GITHUB_REPOSITORY}",
            "${scope}:org.opencontainers.image.documentation=https://github.com/${GITHUB_REPOSITORY}/blob/main/python/README.md",
            "${scope}:org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}",
            "${scope}:org.opencontainers.image.version=${version}",
            "${scope}:org.opencontainers.image.revision=${GITHUB_SHA}",
            "${scope}:org.opencontainers.image.vendor=${GITHUB_ACTOR}",
            "${scope}:org.opencontainers.image.licenses=GPL-3.0-or-later",
            "${scope}:org.opencontainers.image.title=Python ${version}${suffix}",
            "${scope}:org.opencontainers.image.description=Python ${version}${suffix} image.",
            "${scope}:org.opencontainers.image.base.name=${suffix == "-dev" ? "container-base" : "scratch"}"
        ]
    ])
}

# ---------------------------------------------------------------------------
# Cache functions
# ---------------------------------------------------------------------------

function "cache_from" {
    params = [tag]
    result = [
        "type=registry,ref=ghcr.io/${lower(GITHUB_REPOSITORY)}/buildcache:${tag}"
    ]
}

function "cache_to" {
    params = [tag]
    result = [
        "type=registry,ref=ghcr.io/${lower(GITHUB_REPOSITORY)}/buildcache:${tag},compression=zstd,compression-level=12,ignore-error=true"
    ]
}

# ---------------------------------------------------------------------------
# Common settings
# ---------------------------------------------------------------------------

target "common" {
    contexts = {
      "astral-sh/uv" = "docker-image://ghcr.io/astral-sh/uv:latest"
    }
    platforms   = ["linux/arm64", "linux/amd64"]
    output      = ["type=image,compression=zstd,compression-level=19,oci-mediatypes=true,force-compression=true"]
    pull        = true
    attest = [
        "type=provenance,mode=max,builder-id=https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}",
        "type=sbom,generator=docker/buildkit-syft-scanner"
    ]
}

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

group "default" {
    targets = ["apk", "apt"]
}

target "apk" {
    inherits = ["common"]
    name = "python-${version == "latest" ? "latest" : replace(version, ".", "-")}-apk-${variant}"
    matrix = {
        version = PYTHON_VERSIONS_LIST
        variant = ["dev", "minimal"]
    }
    contexts = {
      "container-base" = "docker-image://ghcr.io/wolfi-dev/sdk:latest"
    }
    dockerfile = "apk/Dockerfile"
    target = variant
    args = {
        PYTHON_VERSION = version == "latest" ? "" : version
        STRIP_EXTRAS   = variant == "minimal" ? "true" : "false"
    }
    tags        = python_image_tags(version, variant == "dev" ? "-dev" : "", GITHUB_REPOSITORY, DOCKER_REPOSITORY)
    labels      = oci_labels(version, "-${variant}")
    annotations = oci_annotations(version, "-${variant}")
    cache-from  = cache_from("python-${version == "latest" ? "latest" : replace(version, ".", "-")}-apk-${variant}")
    cache-to    = cache_to("python-${version == "latest" ? "latest" : replace(version, ".", "-")}-apk-${variant}")
}

target "apt" {
    inherits = ["common"]
    name = "python-${version == "latest" ? "latest" : replace(version, ".", "-")}-${distro}-${variant}"
    matrix = {
        version = PYTHON_VERSIONS_LIST
        distro  = ["ubuntu", "debian"]
        variant = ["dev", "minimal"]
    }
    contexts = {
      "container-base" = "docker-image://ghcr.io/${lower(GITHUB_REPOSITORY)}/${distro}:latest"
    }
    dockerfile = "apt/Dockerfile"
    target = variant
    args = {
        PYTHON_VERSION = version == "latest" ? "" : version
        STRIP_EXTRAS   = "false"
    }
    tags        = python_image_tags(version, variant == "dev" ? "-${distro}-dev" : "-${distro}", GITHUB_REPOSITORY, DOCKER_REPOSITORY)
    labels      = oci_labels(version, "-${distro}-${variant}")
    annotations = oci_annotations(version, "-${distro}-${variant}")
    cache-from  = cache_from("python-${version == "latest" ? "latest" : replace(version, ".", "-")}-${distro}-${variant}")
    cache-to    = cache_to("python-${version == "latest" ? "latest" : replace(version, ".", "-")}-${distro}-${variant}")
}
