// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package embedres

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// embedresFS encapsulates the interfaces supported by the embed.FS struct for
// tests.
type embedresFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

// EmbeddedResources knows how to manage an embedded file system of resources
// that are built into the binary.
type EmbeddedResources struct {
	fs embedresFS
}

// NewEmbeddedResources creates a new EmbeddedResources on the given file
// system.
func NewEmbeddedResources(fs embedresFS) *EmbeddedResources {
	return &EmbeddedResources{
		fs: fs,
	}
}

// Bytes returns the bytes for the embedded resource at the given path. If a
// resource at the given path does not exist an error is returned.
func (resources *EmbeddedResources) Bytes(path string) result.Result[[]byte] {
	return result.New(resources.fs.ReadFile(path))
}

// MustBytes returns the bytes for the embedded resource at the given path. If a
// resource at the given path does not exist a panic on the error is raised.
func (resources *EmbeddedResources) MustBytes(path string) []byte {
	return resources.Bytes(path).MustGet()
}

// Paths returns the embedded resource paths that match the given patterns.
// See https://pkg.go.dev/path/filepath#Match for the pattern syntax.
func (resources *EmbeddedResources) Paths(patterns ...string) result.Result[[]string] {
	paths := make([]string, 0, 75)
	fn := resources.matchingPathsFunc(patterns, &paths)
	if err := fs.WalkDir(resources.fs, `.`, fn); err != nil {
		return result.Error[[]string](err)
	}
	sort.Strings(paths)
	return result.Ok(paths)
}

func (resources *EmbeddedResources) isMatch(patterns []string, path string) result.Result[bool] {
	if len(patterns) == 0 {
		return result.Ok(true)
	}
	for _, pattern := range patterns {
		cleanPattern := strings.ReplaceAll(pattern, `\`, `/`)
		match, err := filepath.Match(cleanPattern, path)
		if err != nil {
			return result.Error[bool](err)
		}
		if match {
			return result.Ok(true)
		}
	}
	return result.Ok(false)
}

func (resources *EmbeddedResources) matchingPathsFunc(patterns []string, paths *[]string) fs.WalkDirFunc {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		match := resources.isMatch(patterns, path)
		if match.IsError() {
			return match.Error()
		}
		if match.MustGet() {
			*paths = append(*paths, path)
		}
		return nil
	}
}
