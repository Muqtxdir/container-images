package oci

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// RegistryOptions authenticate with the Docker credentials in DOCKER_CONFIG and stop when ctx is done.
func RegistryOptions(ctx context.Context) []remote.Option {
	return []remote.Option{remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain)}
}

// Push writes index to repository under every tag and prints it as <repository>@<digest>.
func Push(index v1.ImageIndex, repository name.Repository, tags []string, registryOptions ...remote.Option) error {
	indexDigest, err := index.Digest()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		taggedReference, err := name.NewTag(repository.Name() + ":" + tag)
		if err != nil {
			return err
		}
		if err := remote.WriteIndex(taggedReference, index, registryOptions...); err != nil {
			return fmt.Errorf("pushing %s: %w", taggedReference, err)
		}
		log.Printf("pushed %s", taggedReference)
	}
	fmt.Println(repository.Digest(indexDigest.String()))
	return nil
}
