package main

import (
	"slices"
	"testing"

	pep440 "github.com/aquasecurity/go-pep440-version"
	"github.com/google/go-github/v93/github"
)

func TestDistributions(t *testing.T) {
	pythonRelease := release{RepositoryRelease: &github.RepositoryRelease{TagName: "20261026"}}
	for _, assetName := range []string{
		"cpython-3.14.10+20261026-x86_64-unknown-linux-gnu-install_only_stripped.tar.gz",
		"cpython-3.14.10+20261026-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz",
		"cpython-3.14.10+20261026-x86_64-unknown-linux-musl-install_only_stripped.tar.gz",
		"cpython-3.15.2+20261026-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz",
	} {
		pythonRelease.Assets = append(pythonRelease.Assets, &github.ReleaseAsset{Name: github.Ptr(assetName)})
	}

	pythonDistributions, err := pythonRelease.distributions()
	if err != nil {
		t.Fatal(err)
	}
	// musl is ignored, and 3.15.2 is skipped for lacking x86_64.
	if len(pythonDistributions) != 1 || pythonDistributions[0].String() != "3.14.10" {
		t.Errorf("distributions() = %v, want [3.14.10]", pythonDistributions)
	}
}

func TestImageTags(t *testing.T) {
	tests := []struct {
		version  string
		abiFlags string
		want     []string
	}{
		{"3.14.10", "", []string{"3.14.10-debian", "3.14-debian", "debian"}},
		{"3.13.18", "", []string{"3.13.18-debian", "3.13-debian"}},
		{"3.14.10", "t", []string{"3.14.10t-debian", "3.14t-debian"}},
		{"3.16.0a1", "", []string{"3.16.0a1-debian", "3.16-rc-debian"}},
	}
	for _, test := range tests {
		pythonDistribution := distribution{version: pep440.MustParse(test.version), abiFlags: test.abiFlags}
		if got := pythonDistribution.imageTags("debian", "3.14.10"); !slices.Equal(got, test.want) {
			t.Errorf("%s.imageTags() = %q, want %q", pythonDistribution, got, test.want)
		}
	}
}

func TestLatestStableVersion(t *testing.T) {
	var pythonDistributions []distribution
	for _, build := range []struct{ version, abiFlags string }{
		{"3.13.18", ""},
		{"3.14.10", ""},
		{"3.14.10", "t"},
		{"3.16.0a1", ""},
	} {
		pythonDistributions = append(pythonDistributions, distribution{version: pep440.MustParse(build.version), abiFlags: build.abiFlags})
	}
	if got := latestStableVersion(pythonDistributions); got != "3.14.10" {
		t.Errorf("latestStableVersion() = %q, want 3.14.10", got)
	}
}
