// SPDX-FileCopyrightText:  2022, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package processutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sassoftware/sas-ggdk/pkg/result"
)

// executableFunc is a hook for testing.
// nolint: revive
var executableFunc func() (string, error) = os.Executable

// ProcessName returns a simple name for the process. The name is not guaranteed to
// be unique, but it will be the same for subsequent calls of the same process.
func ProcessName() result.Result[string] {
	var name string
	// Ask the OS for name of the running executable; this is a fully qualified
	// filename.
	value, err := executableFunc()
	if err != nil {
		return result.Error[string](err)
	}
	// Discard the directory path, we only want the filename.
	_, name = filepath.Split(value)
	// Trim the filename extension, if present.
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	// Normalize the name by replacing dots with underscores.
	name = strings.ReplaceAll(name, `.`, `_`)
	// Append the process ID.
	pid := os.Getpid()
	name = fmt.Sprintf(`%s-%d`, name, pid)
	return result.Ok(name)
}
