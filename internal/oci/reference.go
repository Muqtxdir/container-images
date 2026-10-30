// Package oci holds what os-images and python-images share: reference
// parsing, the runtime config, labels and building and pushing indexes.
package oci

import (
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// ParseReference parses an image reference, lowercasing its registry and
// repository as registries require; the tag or digest is kept as is.
func ParseReference(value string) (name.Reference, error) {
	nameEnd := strings.LastIndex(value, "/") + 1
	if tagStart := strings.IndexAny(value[nameEnd:], ":@"); tagStart >= 0 {
		nameEnd += tagStart
	} else {
		nameEnd = len(value)
	}
	return name.ParseReference(strings.ToLower(value[:nameEnd]) + value[nameEnd:])
}

// ParseRepository parses a repository, lowercased as registries require.
func ParseRepository(value string) (name.Repository, error) {
	return name.NewRepository(strings.ToLower(value))
}
