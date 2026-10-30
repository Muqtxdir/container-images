package main

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-github/v93/github"

	"github.com/muqtxdir/container-images/internal/oci"
)

var pythonLabels = map[string]string{
	"org.opencontainers.image.title":       "python",
	"org.opencontainers.image.description": "Python is an interpreted, interactive, object-oriented, open-source programming language.",
}

// publisher layers distributions onto one base index and pushes them to repository.
type publisher struct {
	baseIndex     v1.ImageIndex
	baseReference string
	release       release
	origin        oci.Origin
	repository    name.Repository
}

func (p publisher) publish(ctx context.Context, pythonDistribution distribution, imageTags []string) error {
	pythonIndex, err := oci.BuildIndex(p.baseIndex, func(baseImage v1.Image, descriptor v1.Descriptor) (v1.Image, error) {
		releaseAsset, built := pythonDistribution.archives[descriptor.Platform.Architecture]
		if !built {
			return nil, fmt.Errorf("release %s has no python %s for it", p.release.GetTagName(), pythonDistribution)
		}
		pythonLayer, err := buildPythonLayer(ctx, p.release, releaseAsset)
		if err != nil {
			return nil, err
		}
		return p.configurePythonImage(baseImage, descriptor.Digest, pythonDistribution, releaseAsset, pythonLayer)
	})
	if err != nil {
		return err
	}
	return oci.Push(pythonIndex, p.repository, imageTags, oci.RegistryOptions(ctx)...)
}

// configurePythonImage appends pythonLayer to baseImage, keeping the os-images config with PYTHON_VERSION and python3
// as the command.
func (p publisher) configurePythonImage(baseImage v1.Image, baseDigest v1.Hash, pythonDistribution distribution,
	releaseAsset *github.ReleaseAsset, pythonLayer v1.Layer) (v1.Image, error) {
	baseConfigFile, err := baseImage.ConfigFile()
	if err != nil {
		return nil, err
	}
	labels, err := oci.BaseLabels(baseImage, baseDigest, p.baseReference, p.origin)
	if err != nil {
		return nil, err
	}
	releasePublishedAt := p.release.GetPublishedAt().Time
	maps.Copy(labels, pythonLabels)
	labels["org.opencontainers.image.version"] = pythonDistribution.String()
	if releasePublishedAt.After(baseConfigFile.Created.Time) {
		labels["org.opencontainers.image.created"] = releasePublishedAt.UTC().Format(time.RFC3339)
	}

	pythonImage, err := mutate.Append(baseImage, mutate.Addendum{
		Layer: pythonLayer,
		History: v1.History{
			Created:   v1.Time{Time: releasePublishedAt},
			CreatedBy: releaseRepository + " " + p.release.GetTagName() + " " + releaseAsset.GetName(),
		},
	})
	if err != nil {
		return nil, err
	}
	imageConfig := oci.RuntimeConfig
	imageConfig.Env = append(slices.Clone(imageConfig.Env), "PYTHON_VERSION="+pythonDistribution.version.String())
	imageConfig.Cmd = []string{"python3"}
	return oci.Configure(pythonImage, imageConfig, labels)
}
