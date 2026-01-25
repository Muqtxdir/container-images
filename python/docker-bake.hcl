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

variable "PYTHON_LATEST_VERSION" {
    default = "3.14"
}

variable "PYTHON_VERSIONS_LIST" {
    default = [
    "3.8", "3.9", # EOL versions
    "3.10", "3.11", "3.12", "3.13", "3.14"]
}

function "python_image_tags" {
    params = [version, github_repo, docker_repo, latest_version]
    result = concat(
        [
            "ghcr.io/${lower(github_repo)}/python:${version}",
            "docker.io/${lower(docker_repo)}/python:${version}"
        ],
        version == latest_version ? [
            "ghcr.io/${lower(github_repo)}/python:latest",
            "ghcr.io/${lower(github_repo)}/python:3",
            "docker.io/${lower(docker_repo)}/python:latest",
            "docker.io/${lower(docker_repo)}/python:3"
        ] : []
    )
}

target "default" {
    name = "python-${replace(version, ".", "-")}"
    description = "Build Distroless Python images"
    matrix = {
        version = PYTHON_VERSIONS_LIST
    }
    args = {
        USER_GID  = "65532"
        USER_UID  = "65532"
        USER_NAME = "nonroot"
        PYTHON_VERSION = version
    }
    contexts = {
      "container-base" = "docker-image://ghcr.io/wolfi-dev/sdk:latest"
    }
    dockerfile = "Dockerfile"
    platforms = ["linux/amd64", "linux/arm64"]
    output = ["type=image,compression=zstd,compression-level=19,oci-mediatypes=true,force-compression=true"]
    pull = true
    attest = [
        "type=provenance,mode=max",
        "type=sbom,generator=docker/buildkit-syft-scanner"
    ]
    tags = python_image_tags(version, GITHUB_REPOSITORY, DOCKER_REPOSITORY, PYTHON_LATEST_VERSION)
    labels = {
      "org.opencontainers.image.created" = "${timestamp()}"
      "org.opencontainers.image.authors" = "${GITHUB_ACTOR}"
      "org.opencontainers.image.url" = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.documentation" = "https://github.com/${GITHUB_REPOSITORY}/blob/main/python/README.md"
      "org.opencontainers.image.source" = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.version" = "${version}"
      "org.opencontainers.image.revision" = "${GITHUB_SHA}"
      "org.opencontainers.image.vendor" = "${GITHUB_ACTOR}"
      "org.opencontainers.image.licenses" = "GPL-3.0-or-later"
      "org.opencontainers.image.title" = "Distroless Python ${version}"
      "org.opencontainers.image.description" = "Distroless Python ${version} image."
      "org.opencontainers.image.base.name" = "scratch"
    }
    annotations = [
      "index:org.opencontainers.image.created=${timestamp()}",
      "index:org.opencontainers.image.authors=${GITHUB_ACTOR}",
      "index:org.opencontainers.image.url=https://github.com/${GITHUB_REPOSITORY}",
      "index:org.opencontainers.image.documentation=https://github.com/${GITHUB_REPOSITORY}/blob/main/python/README.md",
      "index:org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}",
      "index:org.opencontainers.image.version=${version}",
      "index:org.opencontainers.image.revision=${GITHUB_SHA}",
      "index:org.opencontainers.image.vendor=${GITHUB_ACTOR}",
      "index:org.opencontainers.image.licenses=GPL-3.0-or-later",
      "index:org.opencontainers.image.title=Distroless Python ${version}",
      "index:org.opencontainers.image.description=Distroless Python ${version} image.",
      "index:org.opencontainers.image.base.name=scratch",
      "manifest:org.opencontainers.image.created=${timestamp()}",
      "manifest:org.opencontainers.image.authors=${GITHUB_ACTOR}",
      "manifest:org.opencontainers.image.url=https://github.com/${GITHUB_REPOSITORY}",
      "manifest:org.opencontainers.image.documentation=https://github.com/${GITHUB_REPOSITORY}/blob/main/python/README.md",
      "manifest:org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}",
      "manifest:org.opencontainers.image.version=${version}",
      "manifest:org.opencontainers.image.revision=${GITHUB_SHA}",
      "manifest:org.opencontainers.image.vendor=${GITHUB_ACTOR}",
      "manifest:org.opencontainers.image.licenses=GPL-3.0-or-later",
      "manifest:org.opencontainers.image.title=Distroless Python ${version}",
      "manifest:org.opencontainers.image.description=Distroless Python ${version} image.",
      "manifest:org.opencontainers.image.base.name=scratch"
    ]

}
