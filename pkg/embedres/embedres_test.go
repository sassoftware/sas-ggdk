// SPDX-FileCopyrightText:  2021-2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package embedres_test

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sassoftware/sas-ggdk/pkg/embedres"
	"github.com/sassoftware/sas-ggdk/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata
	dataFS embed.FS
)

func Test_NewEmbeddedResources(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	require.NotNil(t, emb)
	bites := emb.Bytes(`testdata/greeting.txt`)
	require.NoError(t, bites.Error())
	require.NotEmpty(t, t, bites.MustGet())
}

func Test_Paths(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	paths := emb.Paths()
	require.NoError(t, paths.Error())
	require.Len(t, paths.MustGet(), 2)
	require.Equal(t, `testdata/greeting.txt`, paths.MustGet()[0])
	require.Equal(t, `testdata/read.me`, paths.MustGet()[1])
}

func Test_Paths_patternErrBadPattern(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	paths := emb.Paths(`[`)
	require.Error(t, paths.Error())
	require.Same(t, paths.Error(), filepath.ErrBadPattern)
}

func Test_Paths_patternMatch(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	paths := emb.Paths(`testdata/*.txt`)
	require.NoError(t, paths.Error())
	require.Len(t, paths.MustGet(), 1)
	require.Equal(t, `testdata/greeting.txt`, paths.MustGet()[0])
}

func Test_Paths_patternsMatch(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	paths := emb.Paths(`testdata/*.txt`, `testdata/*.me`)
	require.NoError(t, paths.Error())
	require.Len(t, paths.MustGet(), 2)
	require.Equal(t, `testdata/greeting.txt`, paths.MustGet()[0])
	require.Equal(t, `testdata/read.me`, paths.MustGet()[1])
}

func Test_Paths_patternNoMatch(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	paths := emb.Paths(`testdata/*.html`)
	require.NoError(t, paths.Error())
	require.Empty(t, paths.MustGet())
}

func Test_Bytes(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	bites := emb.Bytes(`testdata/greeting.txt`)
	require.NoError(t, bites.Error())
	require.NotEmpty(t, t, bites.MustGet())
	expected := "Hello {{.}}\n"
	actual := string(bites.MustGet())
	actual = strings.ReplaceAll(actual, "\r", "")
	require.Equal(t, expected, actual)
}

func Test_Bytes_failureFileDoesNotExist(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	path := `testdata/hello.txt`
	bites := emb.Bytes(path)
	require.Error(t, bites.Error())
	expected := fmt.Sprintf(`open %s: file does not exist`, path)
	actual := bites.Error().Error()
	require.Equal(t, expected, actual)
}

func Test_MustBytes(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	bites := emb.MustBytes(`testdata/greeting.txt`)
	require.NotEmpty(t, t, bites)
	expected := "Hello {{.}}\n"
	actual := string(bites)
	actual = strings.ReplaceAll(actual, "\r", "")
	require.Equal(t, expected, actual)
}

func Test_MustBytes_failureFileDoesNotExist(t *testing.T) {
	emb := embedres.NewEmbeddedResources(dataFS)
	path := `testdata/hello.txt`
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		actual := err.Error()
		expected := fmt.Sprintf(`open %s: file does not exist`, path)
		require.Equal(t, expected, actual)
	}()
	bites := emb.MustBytes(path)
	require.Empty(t, bites)
}

type readDirFail struct {
	embed.FS
}

func (r *readDirFail) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, errors.New("readdir failed")
}

func Test_FailedDirRead(t *testing.T) {
	fsys := readDirFail{FS: dataFS}
	emb := embedres.NewEmbeddedResources(&fsys)
	paths := emb.Paths(`testdata/greeting.txt`)
	require.ErrorContains(t, paths.Error(), "readdir failed")
}
