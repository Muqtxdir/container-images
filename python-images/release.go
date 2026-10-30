package main

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	pep440 "github.com/aquasecurity/go-pep440-version"
	"github.com/google/go-github/v93/github"
)

const (
	releaseOwner      = "astral-sh"
	releaseRepository = "python-build-standalone"
)

// assetNamePattern matches stripped glibc installs, capturing version, architecture and freethreading, e.g.
// cpython-3.14.10+20261026-x86_64-unknown-linux-gnu-freethreaded-install_only_stripped.tar.gz.
var assetNamePattern = regexp.MustCompile(`^cpython-([^+]+)\+\d+-([^-]+)-unknown-linux-gnu(-freethreaded)?-install_only_stripped\.tar\.gz$`)

// ociArchitectures maps the release architectures we layer to OCI ones.
var ociArchitectures = map[string]string{"x86_64": "amd64", "aarch64": "arm64"}

type release struct {
	*github.RepositoryRelease
	client *github.Client
}

// fetchRelease reads releaseTag, or the latest release, authenticating with GITHUB_TOKEN when it is set.
func fetchRelease(ctx context.Context, releaseTag string) (release, error) {
	var clientOptions []github.ClientOptionsFunc
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		clientOptions = append(clientOptions, github.WithAuthToken(token))
	}
	client, err := github.NewClient(clientOptions...)
	if err != nil {
		return release{}, err
	}

	var repositoryRelease *github.RepositoryRelease
	if releaseTag == "latest" {
		repositoryRelease, _, err = client.Repositories.GetLatestRelease(ctx, releaseOwner, releaseRepository)
	} else {
		repositoryRelease, _, err = client.Repositories.GetReleaseByTag(ctx, releaseOwner, releaseRepository, releaseTag)
	}
	if err != nil {
		return release{}, fmt.Errorf("fetching %s release %s: %w", releaseRepository, releaseTag, err)
	}
	return release{RepositoryRelease: repositoryRelease, client: client}, nil
}

// distributions returns the release's builds that have every architecture, oldest first.
func (r release) distributions() ([]distribution, error) {
	type buildKey struct{ version, abiFlags string }
	archivesByBuild := map[buildKey]map[string]*github.ReleaseAsset{}
	for _, releaseAsset := range r.Assets {
		nameMatch := assetNamePattern.FindStringSubmatch(releaseAsset.GetName())
		if nameMatch == nil {
			continue
		}
		ociArchitecture, layered := ociArchitectures[nameMatch[2]]
		if !layered {
			continue
		}
		key := buildKey{version: nameMatch[1]}
		if nameMatch[3] != "" {
			key.abiFlags = "t"
		}
		if archivesByBuild[key] == nil {
			archivesByBuild[key] = map[string]*github.ReleaseAsset{}
		}
		archivesByBuild[key][ociArchitecture] = releaseAsset
	}

	pythonDistributions := make([]distribution, 0, len(archivesByBuild))
	for key, archivesByArchitecture := range archivesByBuild {
		if len(archivesByArchitecture) < len(ociArchitectures) {
			log.Printf("skipping python %s%s: not built for every architecture", key.version, key.abiFlags)
			continue
		}
		version, err := pep440.Parse(key.version)
		if err != nil {
			return nil, fmt.Errorf("python %s: %w", key.version, err)
		}
		pythonDistributions = append(pythonDistributions, distribution{
			version:  version,
			abiFlags: key.abiFlags,
			archives: archivesByArchitecture,
		})
	}
	if len(pythonDistributions) == 0 {
		return nil, fmt.Errorf("release %s has no CPython builds", r.GetTagName())
	}
	slices.SortFunc(pythonDistributions, func(older, newer distribution) int {
		return cmp.Or(older.version.Compare(newer.version), strings.Compare(older.abiFlags, newer.abiFlags))
	})
	return pythonDistributions, nil
}

func (r release) downloadAsset(ctx context.Context, releaseAsset *github.ReleaseAsset) (io.ReadCloser, error) {
	archiveReader, _, err := r.client.Repositories.DownloadReleaseAsset(ctx, releaseOwner, releaseRepository,
		releaseAsset.GetID(), http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("downloading %s: %w", releaseAsset.GetName(), err)
	}
	return archiveReader, nil
}
