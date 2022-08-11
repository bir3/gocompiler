// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains extra hooks for testing the go command.

//go:build testgo

package work

import (
	"r2.is/gocompiler/src/cmd/gocmd/internal/cfg"
	"r2.is/gocompiler/src/cmd/gocmd/internal/search"
	"fmt"
	       "r2.is/gocompiler/vfs/os"
	"path/filepath"
	"runtime"
)

func init() {
	if v := os.Getenv("TESTGO_VERSION"); v != "" {
		runtimeVersion = v
	}

	if testGOROOT := os.Getenv("TESTGO_GOROOT"); testGOROOT != "" {
		// Disallow installs to the GOROOT from which testgo was built.
		// Installs to other GOROOTs — such as one set explicitly within a test — are ok.
		allowInstall = func(a *Action) error {
			if cfg.BuildN {
				return nil
			}

			rel := search.InDir(a.Target, testGOROOT)
			if rel == "" {
				return nil
			}

			callerPos := ""
			if _, file, line, ok := runtime.Caller(1); ok {
				if shortFile := search.InDir(file, filepath.Join(testGOROOT, "src")); shortFile != "" {
					file = shortFile
				}
				callerPos = fmt.Sprintf("%s:%d: ", file, line)
			}
			return fmt.Errorf("%stestgo must not write to GOROOT (installing to %s)", callerPos, filepath.Join("GOROOT", rel))
		}
	}
}