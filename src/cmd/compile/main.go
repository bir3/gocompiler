// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

import (
	"github.com/bir3/gocompiler/src/cmd/compile/internal/amd64"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/arm"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/arm64"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/base"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/gc"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/loong64"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/mips"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/mips64"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/ppc64"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/riscv64"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/s390x"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/ssagen"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/wasm"
	"github.com/bir3/gocompiler/src/cmd/compile/internal/x86"
	"fmt"
	"github.com/bir3/gocompiler/src/internal/buildcfg"
	"log"
	       "github.com/bir3/gocompiler/vfs/os"
)

var archInits = map[string]func(*ssagen.ArchInfo){
	"386":      x86.Init,
	"amd64":    amd64.Init,
	"arm":      arm.Init,
	"arm64":    arm64.Init,
	"loong64":  loong64.Init,
	"mips":     mips.Init,
	"mipsle":   mips.Init,
	"mips64":   mips64.Init,
	"mips64le": mips64.Init,
	"ppc64":    ppc64.Init,
	"ppc64le":  ppc64.Init,
	"riscv64":  riscv64.Init,
	"s390x":    s390x.Init,
	"wasm":     wasm.Init,
}

func Main() {
	// disable timestamps for reproducible output
	log.SetFlags(0)
	log.SetPrefix("compile: ")

	buildcfg.Check()
	archInit, ok := archInits[buildcfg.GOARCH]
	if !ok {
		fmt.Fprintf(os.Stderr, "compile: unknown architecture %q\n", buildcfg.GOARCH)
		os.Exit(2)
	}

	gc.Main(archInit)
	base.Exit(0)
}
