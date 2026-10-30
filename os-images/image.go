package main

import (
	"fmt"
	"maps"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

var runtimeConfig = v1.Config{
	Env: []string{
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		`PS1=\u@\h:\w\$ `,
		"LANG=C.UTF-8",
	},
	Shell:      []string{"/bin/bash", "-euxo", "pipefail", "-c"},
	StopSignal: "SIGTERM",
	Entrypoint: []string{"tini", "--"},
	Cmd:        []string{"/bin/bash"},
}

func buildIndex(baseIndex v1.ImageIndex, baseReference string, origin buildOrigin) (v1.ImageIndex, error) {
	baseIndexManifest, err := baseIndex.IndexManifest()
	if err != nil {
		return nil, err
	}
	baseIndexDigest, err := baseIndex.Digest()
	if err != nil {
		return nil, err
	}

	index := mutate.IndexMediaType(empty.Index, types.OCIImageIndex)
	var indexAnnotations map[string]string
	for _, descriptor := range baseIndexManifest.Manifests {
		baseImage, err := baseIndex.Image(descriptor.Digest)
		if err != nil {
			return nil, err
		}
		configuredImage, labels, err := configureImage(baseImage, descriptor.Digest, baseReference, origin)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", descriptor.Platform, err)
		}
		index = mutate.AppendManifests(index, mutate.IndexAddendum{
			Add:        configuredImage,
			Descriptor: v1.Descriptor{Platform: descriptor.Platform},
		})
		// Cloned: mutate.Annotations keeps the map it is given.
		indexAnnotations = maps.Clone(labels)
	}
	if indexAnnotations == nil {
		return nil, fmt.Errorf("index %s has no images", baseIndexDigest)
	}
	indexAnnotations["org.opencontainers.image.base.digest"] = baseIndexDigest.String()
	return mutate.Annotations(index, indexAnnotations).(v1.ImageIndex), nil
}

// configureImage applies runtimeConfig, the base's metadata and buildLabels.
func configureImage(baseImage v1.Image, baseDigest v1.Hash, baseReference string, origin buildOrigin) (v1.Image, map[string]string, error) {
	baseConfigFile, err := baseImage.ConfigFile()
	if err != nil {
		return nil, nil, err
	}
	baseManifest, err := baseImage.Manifest()
	if err != nil {
		return nil, nil, err
	}

	labels := maps.Clone(baseConfigFile.Config.Labels)
	if labels == nil {
		labels = map[string]string{}
	}
	maps.Copy(labels, baseManifest.Annotations)
	maps.Copy(labels, buildLabels(baseReference, baseDigest, baseConfigFile.Created.Time.UTC(), origin))

	config := runtimeConfig
	config.Labels = labels
	configuredImage, err := mutate.Config(baseImage, config)
	if err != nil {
		return nil, nil, err
	}
	return mutate.Annotations(configuredImage, labels).(v1.Image), labels, nil
}
