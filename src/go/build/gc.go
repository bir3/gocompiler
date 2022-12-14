// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc

package build

import (
	"path/filepath"
	"runtime"
"github.com/bir3/gocompiler/vfs"
)

// getToolDir returns the default value of ToolDir.
func getToolDir() string {
	return filepath.Join(vfs.GOROOT, "pkg/tool/"+runtime.GOOS+"_"+runtime.GOARCH)
}
