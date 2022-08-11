// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s390x

import (
	"r2.is/gocompiler/src/cmd/compile/internal/ssagen"
	"r2.is/gocompiler/src/cmd/internal/obj/s390x"
)

func Init(arch *ssagen.ArchInfo) {
	arch.LinkArch = &s390x.Links390x
	arch.REGSP = s390x.REGSP
	arch.MAXWIDTH = 1 << 50

	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = ssaMarkMoves
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
}
