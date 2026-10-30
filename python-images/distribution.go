package main

import (
	"slices"
	"strings"

	pep440 "github.com/aquasecurity/go-pep440-version"
	"github.com/google/go-github/v93/github"
	"golang.org/x/mod/semver"
)

// distribution is one CPython build in a release.
type distribution struct {
	version  pep440.Version
	abiFlags string                          // "t" when freethreaded, as in sys.abiflags.
	archives map[string]*github.ReleaseAsset // By OCI architecture.
}

// String returns the version with its ABI flags, e.g. 3.14.10t.
func (d distribution) String() string {
	return d.version.String() + d.abiFlags
}

func (d distribution) isStable() bool {
	return !d.version.IsPreRelease() && d.abiFlags == ""
}

// imageTags returns e.g. 3.14.10t-<variant> and 3.14t-<variant>, 3.16-rc-<variant> for pre-releases, and also
// <variant> for latestVersion.
func (d distribution) imageTags(variant, latestVersion string) []string {
	minorTag := strings.TrimPrefix(semver.MajorMinor("v"+d.version.BaseVersion()), "v") + d.abiFlags
	if d.version.IsPreRelease() {
		minorTag += "-rc"
	}
	imageTags := []string{d.String() + "-" + variant, minorTag + "-" + variant}
	if d.String() == latestVersion {
		imageTags = append(imageTags, variant)
	}
	return imageTags
}

// latestStableVersion returns the newest stable version in sorted pythonDistributions, or "" when there is none.
func latestStableVersion(pythonDistributions []distribution) string {
	for _, pythonDistribution := range slices.Backward(pythonDistributions) {
		if pythonDistribution.isStable() {
			return pythonDistribution.String()
		}
	}
	return ""
}
