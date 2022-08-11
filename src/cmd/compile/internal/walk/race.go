// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"r2.is/gocompiler/src/cmd/compile/internal/base"
	"r2.is/gocompiler/src/cmd/compile/internal/ir"
	"r2.is/gocompiler/src/cmd/compile/internal/types"
	"r2.is/gocompiler/src/cmd/internal/src"
)

func instrument(fn *ir.Func) {
	if fn.Pragma&ir.Norace != 0 || (fn.Linksym() != nil && fn.Linksym().ABIWrapper()) {
		return
	}

	if !base.Flag.Race || !base.Compiling(base.NoRacePkgs) {
		fn.SetInstrumentBody(true)
	}

	if base.Flag.Race {
		lno := base.Pos
		base.Pos = src.NoXPos
		var init ir.Nodes
		fn.Enter.Prepend(mkcallstmt("racefuncenter", mkcall("getcallerpc", types.Types[types.TUINTPTR], &init)))
		if len(init) != 0 {
			base.Fatalf("race walk: unexpected init for getcallerpc")
		}
		fn.Exit.Append(mkcallstmt("racefuncexit"))
		base.Pos = lno
	}
}
