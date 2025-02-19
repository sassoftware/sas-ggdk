// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package filters_test

import (
	"io/fs"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/filters"
	"github.com/stretchr/testify/require"
)

func Test_NewIsDirectoryFilter(t *testing.T) {
	filter := filters.NewIsDirectoryFilter()
	dir := newFakeDirEntry(`abc`, true)
	state := filter(dir)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
	file := newFakeDirEntry(`def.txt`, false)
	state = filter(file)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
}

func Test_NewIsFileFilter(t *testing.T) {
	filter := filters.NewIsFileFilter()
	dir := newFakeDirEntry(`abc`, true)
	state := filter(dir)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
	file := newFakeDirEntry(`def.txt`, false)
	state = filter(file)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
}

func Test_NewIsEmptyMapFilter_slice(t *testing.T) {
	filter := filters.NewIsEmptyMapFilter[map[string]string, string, string]()
	m := make(map[string]string, 1)
	state := filter(m)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
	m[`color`] = `red`
	state = filter(m)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
	state = filter(nil)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
}

func Test_NewIsEmptyStringFilter_string(t *testing.T) {
	filter := filters.NewIsEmptyStringFilter()
	state := filter(`hello`)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
	state = filter(``)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
}

func Test_NewIsEqualFilter(t *testing.T) {
	filter := filters.NewIsEqualFilter(10)
	state := filter(10)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
	state = filter(20)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
}

func Test_NewIsEmptySliceFilter_slice(t *testing.T) {
	filter := filters.NewIsEmptySliceFilter[string]()
	state := filter([]string{`hello`})
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
	state = filter([]string{})
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
	state = filter(nil)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
}

func Test_NewMapContainsKeyFilter(t *testing.T) {
	m := make(map[string]string)
	filter := filters.NewMapContainsKeyFilter(m)
	state := filter(`color`)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
	m[`color`] = `red`
	state = filter(`color`)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
}

func Test_NewSliceContainsFilter(t *testing.T) {
	colors := []string{`red`, `blue`, `green`}
	filter := filters.NewSliceContainsFilter(colors)
	state := filter(`yellow`)
	require.NoError(t, state.Error())
	require.False(t, state.MustGet())
	state = filter(`blue`)
	require.NoError(t, state.Error())
	require.True(t, state.MustGet())
}

// newFakeDirEntry returns a new fake dir entry for testing.
func newFakeDirEntry(name string, isDir bool) *fakeDirEntry {
	return &fakeDirEntry{
		name:  name,
		isDir: isDir,
	}
}

// fakeDirEntry represents a fake dir entry that is either a regular
// file or a directory. It has a name, type, and info.
type fakeDirEntry struct {
	name  string
	isDir bool
}

// Name returns the name of the fake dir entry.
func (f fakeDirEntry) Name() string {
	return f.name
}

// IsDir returns true if the fake dir entry is a dir, false otherwise.
func (f fakeDirEntry) IsDir() bool {
	return f.isDir
}

// Type returns the type of the fake dir entry.
func (f fakeDirEntry) Type() fs.FileMode {
	panic(`implement me`)
}

// Info returns the info of the fake dir entry.
func (f fakeDirEntry) Info() (fs.FileInfo, error) {
	panic(`implement me`)
}
