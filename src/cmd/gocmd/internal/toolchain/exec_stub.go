// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js || wasip1

package toolchain

import "github.com/bir3/gocompiler/src/cmd/gocmd/internal/base"

func execGoToolchain(gotoolchain, dir, exe string) {
	base.Fatalf("execGoToolchain unsupported")
}
