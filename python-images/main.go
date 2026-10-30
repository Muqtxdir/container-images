package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"golang.org/x/sync/errgroup"

	"github.com/muqtxdir/container-images/internal/oci"
)

// maxConcurrentBuilds is where more builds stop being faster and only cost memory.
const maxConcurrentBuilds = 4

type options struct {
	source     name.Reference
	repository name.Repository
	variant    string
	releaseTag string
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
	pythonRelease, err := fetchRelease(ctx, parsedOptions.releaseTag)
	if err != nil {
		return err
	}
	pythonDistributions, err := pythonRelease.distributions()
	if err != nil {
		return err
	}

	imagePublisher := publisher{
		baseIndex:     baseIndex,
		baseReference: parsedOptions.source.String(),
		release:       pythonRelease,
		origin:        oci.OriginFromEnv(),
		repository:    parsedOptions.repository,
	}
	latestVersion := latestStableVersion(pythonDistributions)
	buildErrors := make([]error, len(pythonDistributions))
	var builds errgroup.Group
	builds.SetLimit(maxConcurrentBuilds)
	for position, pythonDistribution := range pythonDistributions {
		builds.Go(func() error {
			imageTags := pythonDistribution.imageTags(parsedOptions.variant, latestVersion)
			if err := imagePublisher.publish(ctx, pythonDistribution, imageTags); err != nil {
				buildErrors[position] = fmt.Errorf("python %s: %w", pythonDistribution, err)
			}
			return nil
		})
	}
	builds.Wait()
	return errors.Join(buildErrors...)
}

func parseOptions() (options, error) {
	parsedOptions := options{releaseTag: "latest"}
	flag.Func("source", "multi-platform base image", func(value string) (err error) {
		parsedOptions.source, err = oci.ParseReference(value)
		return err
	})
	flag.Func("repository", "repository to push to", func(value string) (err error) {
		parsedOptions.repository, err = oci.ParseRepository(value)
		return err
	})
	flag.StringVar(&parsedOptions.variant, "variant", "", "tag suffix")
	flag.StringVar(&parsedOptions.releaseTag, "release", parsedOptions.releaseTag, "python-build-standalone release")
	flag.Parse()
	if parsedOptions.source == nil || parsedOptions.repository.Name() == "" || parsedOptions.variant == "" {
		return parsedOptions, errors.New("--source, --repository and --variant are required")
	}
	return parsedOptions, nil
}
