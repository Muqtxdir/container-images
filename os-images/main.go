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

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/muqtxdir/container-images/internal/oci"
)

type options struct {
	source     name.Reference
	repository name.Repository
	tags       []string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	parsedOptions, err := parseOptions()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	baseIndex, err := remote.Index(parsedOptions.source, oci.RegistryOptions(ctx)...)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", parsedOptions.source, err)
	}
	buildOrigin := oci.OriginFromEnv()
	configuredIndex, err := oci.BuildIndex(baseIndex, func(baseImage v1.Image, descriptor v1.Descriptor) (v1.Image, error) {
		labels, err := oci.BaseLabels(baseImage, descriptor.Digest, parsedOptions.source.String(), buildOrigin)
		if err != nil {
			return nil, err
		}
		return oci.Configure(baseImage, oci.RuntimeConfig, labels)
	})
	if err != nil {
		return fmt.Errorf("building %s: %w", parsedOptions.source, err)
	}
	return oci.Push(configuredIndex, parsedOptions.repository, parsedOptions.tags, oci.RegistryOptions(ctx)...)
}

func parseOptions() (options, error) {
	var parsedOptions options
	flag.Func("source", "multi-platform base image", func(value string) (err error) {
		parsedOptions.source, err = oci.ParseReference(value)
		return err
	})
	flag.Func("repository", "repository to push to", func(value string) (err error) {
		parsedOptions.repository, err = oci.ParseRepository(value)
		return err
	})
	flag.Func("tags", "tags to push, separated by spaces", func(value string) error {
		parsedOptions.tags = strings.Fields(value)
		return nil
	})
	flag.Parse()
	if parsedOptions.source == nil || parsedOptions.repository.Name() == "" || len(parsedOptions.tags) == 0 {
		return parsedOptions, errors.New("--source, --repository and --tags are required")
	}
	return parsedOptions, nil
}
