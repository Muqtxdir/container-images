package oci

import (
	"maps"
	"os"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Origin is where the image is built.
type Origin struct {
	repositoryURL string
	vendor        string
	authors       string
	revision      string
}

// Labels returns the OCI labels for this build, omitting empty values.
func Labels(baseReference string, baseDigest v1.Hash, created time.Time, origin Origin) map[string]string {
	labels := map[string]string{
		"org.opencontainers.image.created":     created.Format(time.RFC3339),
		"org.opencontainers.image.url":         origin.repositoryURL,
		"org.opencontainers.image.source":      origin.repositoryURL,
		"org.opencontainers.image.revision":    origin.revision,
		"org.opencontainers.image.vendor":      origin.vendor,
		"org.opencontainers.image.authors":     origin.authors,
		"org.opencontainers.image.base.name":   baseReference,
		"org.opencontainers.image.base.digest": baseDigest.String(),
	}
	maps.DeleteFunc(labels, func(_, value string) bool { return value == "" })
	return labels
}

// OriginFromEnv reads the GitHub Actions variables.
func OriginFromEnv() Origin {
	serverURL := os.Getenv("GITHUB_SERVER_URL")
	repository := os.Getenv("GITHUB_REPOSITORY")
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	origin := Origin{vendor: owner, revision: os.Getenv("GITHUB_SHA")}
	if serverURL != "" && repository != "" {
		origin.repositoryURL = serverURL + "/" + repository
	}
	if serverURL != "" && owner != "" {
		origin.authors = serverURL + "/" + owner
	}
	return origin
}
