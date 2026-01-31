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

variable "NODE_LATEST_VERSION" {
    default = "25"
}

variable "NODE_LTS_VERSION" {
    default = "24"
}

variable "NODE_VERSIONS_LIST" {
    default = [
    "16", "18", "21", "23", # EOL versions
    "20", "22", "24", # LTS versions
    "25"] # Current version
}

function "node_image_tags" {
    params = [version, github_repo, docker_repo, latest_version, lts_version]
    result = concat(
        [
            "ghcr.io/${lower(github_repo)}/node:${version}",
            "docker.io/${lower(docker_repo)}/node:${version}"
        ],
        version == latest_version ? [
            "ghcr.io/${lower(github_repo)}/node:latest",
            "docker.io/${lower(docker_repo)}/node:latest"
        ] : [],
        version == lts_version ? [
            "ghcr.io/${lower(github_repo)}/node:lts",
            "docker.io/${lower(docker_repo)}/node:lts"
        ] : []
    )
}

target "default" {
    name = "node-${version}"
    description = "Build Distroless Node.js images"
    matrix = {
        version = NODE_VERSIONS_LIST
    }
    args = {
        USER_GID  = "65532"
        USER_UID  = "65532"
        USER_NAME = "nonroot"
        NODE_VERSION = version
    }
    contexts = {
      "container-base" = "docker-image://ghcr.io/wolfi-dev/sdk:latest"
    }
    dockerfile = "Dockerfile"
    platforms = ["linux/amd64", "linux/arm64"]
    output = ["type=image,compression=zstd,compression-level=19,oci-mediatypes=true,force-compression=true"]
    pull = true
    attest = [
        "type=provenance,mode=max,builder-id=https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}",
        "type=sbom,generator=docker/buildkit-syft-scanner"
    ]
    tags = node_image_tags(version, GITHUB_REPOSITORY, DOCKER_REPOSITORY, NODE_LATEST_VERSION, NODE_LTS_VERSION)
    labels = {
      "org.opencontainers.image.created" = "${timestamp()}"
      "org.opencontainers.image.authors" = "${GITHUB_ACTOR}"
      "org.opencontainers.image.url" = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.documentation" = "https://github.com/${GITHUB_REPOSITORY}/blob/main/node/README.md"
      "org.opencontainers.image.source" = "https://github.com/${GITHUB_REPOSITORY}"
      "org.opencontainers.image.version" = "${version}"
      "org.opencontainers.image.revision" = "${GITHUB_SHA}"
      "org.opencontainers.image.vendor" = "${GITHUB_ACTOR}"
      "org.opencontainers.image.licenses" = "GPL-3.0-or-later"
      "org.opencontainers.image.title" = "Distroless Node.js ${version}"
      "org.opencontainers.image.description" = "Distroless Node.js ${version} image."
      "org.opencontainers.image.base.name" = "scratch"
    }
    annotations = [
      "index:org.opencontainers.image.created=${timestamp()}",
      "index:org.opencontainers.image.authors=${GITHUB_ACTOR}",
      "index:org.opencontainers.image.url=https://github.com/${GITHUB_REPOSITORY}",
      "index:org.opencontainers.image.documentation=https://github.com/${GITHUB_REPOSITORY}/blob/main/node/README.md",
      "index:org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}",
      "index:org.opencontainers.image.version=${version}",
      "index:org.opencontainers.image.revision=${GITHUB_SHA}",
      "index:org.opencontainers.image.vendor=${GITHUB_ACTOR}",
      "index:org.opencontainers.image.licenses=GPL-3.0-or-later",
      "index:org.opencontainers.image.title=Distroless Node.js ${version}",
      "index:org.opencontainers.image.description=Distroless Node.js ${version} image.",
      "index:org.opencontainers.image.base.name=scratch",
      "manifest:org.opencontainers.image.created=${timestamp()}",
      "manifest:org.opencontainers.image.authors=${GITHUB_ACTOR}",
      "manifest:org.opencontainers.image.url=https://github.com/${GITHUB_REPOSITORY}",
      "manifest:org.opencontainers.image.documentation=https://github.com/${GITHUB_REPOSITORY}/blob/main/node/README.md",
      "manifest:org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}",
      "manifest:org.opencontainers.image.version=${version}",
      "manifest:org.opencontainers.image.revision=${GITHUB_SHA}",
      "manifest:org.opencontainers.image.vendor=${GITHUB_ACTOR}",
      "manifest:org.opencontainers.image.licenses=GPL-3.0-or-later",
      "manifest:org.opencontainers.image.title=Distroless Node.js ${version}",
      "manifest:org.opencontainers.image.description=Distroless Node.js ${version} image.",
      "manifest:org.opencontainers.image.base.name=scratch"
    ]

}
