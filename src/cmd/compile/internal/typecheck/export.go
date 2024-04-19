// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typecheck

import (
	"github.com/bir3/gocompiler/src/cmd/compile/internal/base"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/ir"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/types"
	"github.com/bir3/gocompiler/src/cmd/internal/src"
)

// importfunc declares symbol s as an imported function with type t.
func importfunc(s *types.Sym, t *types.Type) {
	fn := ir.NewFunc(src.NoXPos, src.NoXPos, s, t)
	importsym(fn.Nname)
}

// importvar declares symbol s as an imported variable with type t.
func importvar(s *types.Sym, t *types.Type) {
	n := ir.NewNameAt(src.NoXPos, s, t)
	n.Class = ir.PEXTERN
	importsym(n)
}

func importsym(name *ir.Name) {
	sym := name.Sym()
	if sym.Def != nil {
		base.Fatalf("importsym of symbol that already exists: %v", sym.Def)
	}
	sym.Def = name
}
