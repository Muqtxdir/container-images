package oci

import (
	"fmt"
	"maps"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// RuntimeConfig is the config os-images gives every image.
var RuntimeConfig = v1.Config{
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

// BuildIndex returns an OCI index of configure applied to each image of baseIndex, annotated with their labels.
func BuildIndex(baseIndex v1.ImageIndex, configure func(baseImage v1.Image, descriptor v1.Descriptor) (v1.Image, error)) (v1.ImageIndex, error) {
	baseIndexManifest, err := baseIndex.IndexManifest()
	if err != nil {
		return nil, err
	}
	baseIndexDigest, err := baseIndex.Digest()
	if err != nil {
		return nil, err
	}

	index := mutate.IndexMediaType(empty.Index, types.OCIImageIndex)
	var annotations map[string]string
	for _, descriptor := range baseIndexManifest.Manifests {
		baseImage, err := baseIndex.Image(descriptor.Digest)
		if err != nil {
			return nil, err
		}
		image, err := configure(baseImage, descriptor)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", descriptor.Platform, err)
		}
		configFile, err := image.ConfigFile()
		if err != nil {
			return nil, err
		}
		index = mutate.AppendManifests(index, mutate.IndexAddendum{
			Add:        image,
			Descriptor: v1.Descriptor{Platform: descriptor.Platform},
		})
		// mutate.Annotations keeps the map it is given.
		annotations = maps.Clone(configFile.Config.Labels)
	}
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["org.opencontainers.image.base.digest"] = baseIndexDigest.String()
	return mutate.Annotations(index, annotations).(v1.ImageIndex), nil
}

// BaseLabels returns baseImage's labels and annotations overlaid with this build's Labels.
func BaseLabels(baseImage v1.Image, baseDigest v1.Hash, baseReference string, origin Origin) (map[string]string, error) {
	baseConfigFile, err := baseImage.ConfigFile()
	if err != nil {
		return nil, err
	}
	baseManifest, err := baseImage.Manifest()
	if err != nil {
		return nil, err
	}

	labels := map[string]string{}
	maps.Copy(labels, baseConfigFile.Config.Labels)
	maps.Copy(labels, baseManifest.Annotations)
	maps.Copy(labels, Labels(baseReference, baseDigest, baseConfigFile.Created.Time.UTC(), origin))
	return labels, nil
}

// Configure sets config with labels on image and annotates it with them.
func Configure(image v1.Image, config v1.Config, labels map[string]string) (v1.Image, error) {
	config.Labels = labels
	configured, err := mutate.Config(image, config)
	if err != nil {
		return nil, err
	}
	return mutate.Annotations(configured, labels).(v1.Image), nil
}
