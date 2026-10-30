// Command os-images configures and labels a container-fs image, pushes it to every --target and prints each as <repository>@<digest>.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var source name.Reference
	var targets []name.Reference
	flag.Func("source", "multi-platform base image, e.g. ghcr.io/muqtxdir/container-fs/debian:stable", func(value string) (err error) {
		source, err = parseReference(value)
		return err
	})
	flag.Func("target", "image to push, e.g. ghcr.io/muqtxdir/container-images/debian:stable (repeatable)", func(value string) error {
		targetReference, err := parseReference(value)
		targets = append(targets, targetReference)
		return err
	})
	flag.Parse()
	if source == nil || len(targets) == 0 {
		return errors.New("--source and at least one --target are required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	registryOptions := []remote.Option{remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain)}

	baseIndex, err := remote.Index(source, registryOptions...)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", source, err)
	}
	index, err := buildIndex(baseIndex, source.String(), buildOriginFromEnv())
	if err != nil {
		return fmt.Errorf("building %s: %w", source, err)
	}
	indexDigest, err := index.Digest()
	if err != nil {
		return err
	}
	printedRepositories := map[name.Repository]bool{}
	for _, target := range targets {
		if err := remote.WriteIndex(target, index, registryOptions...); err != nil {
			return fmt.Errorf("pushing %s: %w", target, err)
		}
		log.Printf("pushed %s", target)
		if repository := target.Context(); !printedRepositories[repository] {
			printedRepositories[repository] = true
			fmt.Println(repository.Digest(indexDigest.String()))
		}
	}
	return nil
}

// parseReference parses an image reference, lowercasing its registry and
// repository as registries require; the tag or digest is kept as is.
func parseReference(value string) (name.Reference, error) {
	nameEnd := strings.LastIndex(value, "/") + 1
	if tagStart := strings.IndexAny(value[nameEnd:], ":@"); tagStart >= 0 {
		nameEnd += tagStart
	} else {
		nameEnd = len(value)
	}
	return name.ParseReference(strings.ToLower(value[:nameEnd]) + value[nameEnd:])
}
