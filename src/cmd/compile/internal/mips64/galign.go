// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mips64

import (
	"github.com/bir3/gocompiler/src/cmd/compile/internal/ssa"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/ssagen"
	"github.com/bir3/gocompiler/src/cmd/internal/obj/mips"
	"github.com/bir3/gocompiler/src/internal/buildcfg"
)

func Init(arch *ssagen.ArchInfo) {
	arch.LinkArch = &mips.Linkmips64
	if buildcfg.GOARCH == "mips64le" {
		arch.LinkArch = &mips.Linkmips64le
	}
	arch.REGSP = mips.REGSP
	arch.MAXWIDTH = 1 << 50
	arch.SoftFloat = buildcfg.GOMIPS64 == "softfloat"
	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = func(s *ssagen.State, b *ssa.Block) {}
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
}
