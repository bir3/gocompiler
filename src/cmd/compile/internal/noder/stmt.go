// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package noder

import (
	"github.com/bir3/gocompiler/src/cmd/compile/internal/ir"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/syntax"
)

// TODO(mdempsky): Investigate replacing with switch statements or dense arrays.

var branchOps = [...]ir.Op{
	syntax.Break:		ir.OBREAK,
	syntax.Continue:	ir.OCONTINUE,
	syntax.Fallthrough:	ir.OFALL,
	syntax.Goto:		ir.OGOTO,
}

var callOps = [...]ir.Op{
	syntax.Defer:	ir.ODEFER,
	syntax.Go:	ir.OGO,
}
