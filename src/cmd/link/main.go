// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"r2.is/gocompiler/src/cmd/internal/sys"
	"r2.is/gocompiler/src/cmd/link/internal/amd64"
	"r2.is/gocompiler/src/cmd/link/internal/arm"
	"r2.is/gocompiler/src/cmd/link/internal/arm64"
	"r2.is/gocompiler/src/cmd/link/internal/ld"
	"r2.is/gocompiler/src/cmd/link/internal/mips"
	"r2.is/gocompiler/src/cmd/link/internal/mips64"
	"r2.is/gocompiler/src/cmd/link/internal/ppc64"
	"r2.is/gocompiler/src/cmd/link/internal/riscv64"
	"r2.is/gocompiler/src/cmd/link/internal/s390x"
	"r2.is/gocompiler/src/cmd/link/internal/wasm"
	"r2.is/gocompiler/src/cmd/link/internal/x86"
	"fmt"
	"r2.is/gocompiler/src/internal/buildcfg"
	       "r2.is/gocompiler/vfs/os"
)

// The bulk of the linker implementation lives in cmd/link/internal/ld.
// Architecture-specific code lives in cmd/link/internal/GOARCH.
//
// Program initialization:
//
// Before any argument parsing is done, the Init function of relevant
// architecture package is called. The only job done in Init is
// configuration of the architecture-specific variables.
//
// Then control flow passes to ld.Main, which parses flags, makes
// some configuration decisions, and then gives the architecture
// packages a second chance to modify the linker's configuration
// via the ld.Arch.Archinit function.

func Main() {
	var arch *sys.Arch
	var theArch ld.Arch

	buildcfg.Check()
	switch buildcfg.GOARCH {
	default:
		fmt.Fprintf(os.Stderr, "link: unknown architecture %q\n", buildcfg.GOARCH)
		os.Exit(2)
	case "386":
		arch, theArch = x86.Init()
	case "amd64":
		arch, theArch = amd64.Init()
	case "arm":
		arch, theArch = arm.Init()
	case "arm64":
		arch, theArch = arm64.Init()
	case "mips", "mipsle":
		arch, theArch = mips.Init()
	case "mips64", "mips64le":
		arch, theArch = mips64.Init()
	case "ppc64", "ppc64le":
		arch, theArch = ppc64.Init()
	case "riscv64":
		arch, theArch = riscv64.Init()
	case "s390x":
		arch, theArch = s390x.Init()
	case "wasm":
		arch, theArch = wasm.Init()
	}
	ld.Main(arch, theArch)
}
