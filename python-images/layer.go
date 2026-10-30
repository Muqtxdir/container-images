package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/google/go-github/v93/github"
	"github.com/klauspost/compress/zstd"
	"github.com/opencontainers/go-digest"
)

const (
	archivePrefix = "python/"
	installPrefix = "usr/local/"
)

// excludedPaths leave out pip, libpython (python is statically linked), Tcl/Tk and its users, and the terminfo and
// man pages the bases provide.
var excludedPaths = []string{
	"bin/idle*",
	"bin/pip*",
	"lib/itcl*/**",
	"lib/libpython*",
	"lib/libtcl*",
	"lib/tcl*/**",
	"lib/thread*/**",
	"lib/tk*/**",
	"lib/python3.*/ensurepip/**",
	"lib/python3.*/idlelib/**",
	"lib/python3.*/lib-dynload/_tkinter*",
	"lib/python3.*/site-packages/pip/**",
	"lib/python3.*/site-packages/pip-*/**",
	"lib/python3.*/tkinter/**",
	"lib/python3.*/turtle.py",
	"lib/python3.*/turtledemo/**",
	"share/man/**",
	"share/terminfo/**",
}

// buildPythonLayer streams releaseAsset through digest verification and repacking into a zstd layer held in memory.
func buildPythonLayer(ctx context.Context, pythonRelease release, releaseAsset *github.ReleaseAsset) (v1.Layer, error) {
	expectedDigest, err := digest.Parse(releaseAsset.GetDigest())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", releaseAsset.GetName(), err)
	}
	archiveReader, err := pythonRelease.downloadAsset(ctx, releaseAsset)
	if err != nil {
		return nil, err
	}
	defer archiveReader.Close()

	digestVerifier := expectedDigest.Verifier()
	var compressedLayer bytes.Buffer
	if err := repackArchive(io.TeeReader(archiveReader, digestVerifier), &compressedLayer); err != nil {
		return nil, fmt.Errorf("repacking %s: %w", releaseAsset.GetName(), err)
	}
	// gzip stops at its trailer; hash anything after it too.
	if _, err := io.Copy(digestVerifier, archiveReader); err != nil {
		return nil, fmt.Errorf("downloading %s: %w", releaseAsset.GetName(), err)
	}
	if !digestVerifier.Verified() {
		return nil, fmt.Errorf("%s does not match %s", releaseAsset.GetName(), expectedDigest)
	}
	return tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(compressedLayer.Bytes())), nil
	}, tarball.WithMediaType(types.OCILayerZStd))
}

// repackArchive writes the archive's kept entries to layerWriter as a zstd tar under installPrefix, owned by root and
// not writable by group or others.
func repackArchive(archiveReader io.Reader, layerWriter io.Writer) error {
	gzipReader, err := gzip.NewReader(archiveReader)
	if err != nil {
		return err
	}
	// Builds already run concurrently, so one encoder goroutine each is as fast and far lighter.
	zstdWriter, err := zstd.NewWriter(layerWriter, zstd.WithEncoderLevel(zstd.SpeedBestCompression), zstd.WithEncoderConcurrency(1))
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzipReader)
	tarWriter := tar.NewWriter(zstdWriter)
	for {
		entryHeader, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		relativePath, insideArchive := strings.CutPrefix(entryHeader.Name, archivePrefix)
		if !insideArchive {
			return fmt.Errorf("%s is outside %s", entryHeader.Name, archivePrefix)
		}
		if isExcluded(relativePath) {
			continue
		}

		entryHeader.Name = installPrefix + relativePath
		if entryHeader.Typeflag == tar.TypeLink {
			entryHeader.Linkname = installPrefix + strings.TrimPrefix(entryHeader.Linkname, archivePrefix)
		}
		entryHeader.Mode &^= 0o022
		entryHeader.Uid, entryHeader.Gid, entryHeader.Uname, entryHeader.Gname = 0, 0, "", ""
		if err := tarWriter.WriteHeader(entryHeader); err != nil {
			return err
		}
		if _, err := io.Copy(tarWriter, tarReader); err != nil {
			return err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return err
	}
	return zstdWriter.Close()
}

func isExcluded(relativePath string) bool {
	return slices.ContainsFunc(excludedPaths, func(excludedPattern string) bool {
		return doublestar.MatchUnvalidated(excludedPattern, relativePath)
	})
}
