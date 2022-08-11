// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"bufio"
	"github.com/bir3/gocompiler/src/cmd/asm/flag"
	"fmt"
	"github.com/bir3/gocompiler/src/internal/buildcfg"
	"log"
	       "github.com/bir3/gocompiler/vfs/os"

	"github.com/bir3/gocompiler/src/cmd/asm/internal/arch"
	"github.com/bir3/gocompiler/src/cmd/asm/internal/asm"
	"github.com/bir3/gocompiler/src/cmd/asm/internal/flags"
	"github.com/bir3/gocompiler/src/cmd/asm/internal/lex"

	"github.com/bir3/gocompiler/src/cmd/internal/bio"
	"github.com/bir3/gocompiler/src/cmd/internal/obj"
	"github.com/bir3/gocompiler/src/cmd/internal/objabi"
)

func Main() {
	log.SetFlags(0)
	log.SetPrefix("asm: ")

	buildcfg.Check()
	GOARCH := buildcfg.GOARCH

	flags.Parse()

	architecture := arch.Set(GOARCH, *flags.Shared || *flags.Dynlink)
	if architecture == nil {
		log.Fatalf("unrecognized architecture %s", GOARCH)
	}

	ctxt := obj.Linknew(architecture.LinkArch)
	ctxt.Debugasm = flags.PrintOut
	ctxt.Debugvlog = flags.DebugV
	ctxt.Flag_dynlink = *flags.Dynlink
	ctxt.Flag_linkshared = *flags.Linkshared
	ctxt.Flag_shared = *flags.Shared || *flags.Dynlink
	ctxt.Flag_maymorestack = flags.DebugFlags.MayMoreStack
	ctxt.IsAsm = true
	ctxt.Pkgpath = *flags.Importpath
	switch *flags.Spectre {
	default:
		log.Printf("unknown setting -spectre=%s", *flags.Spectre)
		os.Exit(2)
	case "":
		// nothing
	case "index":
		// known to compiler; ignore here so people can use
		// the same list with -gcflags=-spectre=LIST and -asmflags=-spectrre=LIST
	case "all", "ret":
		ctxt.Retpoline = true
	}

	ctxt.Bso = bufio.NewWriter(os.Stdout)
	defer ctxt.Bso.Flush()

	architecture.Init(ctxt)

	// Create object file, write header.
	buf, err := bio.Create(*flags.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer buf.Close()

	if !*flags.SymABIs {
		buf.WriteString(objabi.HeaderString())
		fmt.Fprintf(buf, "!\n")
	}

	var ok, diag bool
	var failedFile string
	for _, f := range flag.Args() {
		lexer := lex.NewLexer(f)
		parser := asm.NewParser(ctxt, architecture, lexer,
			*flags.CompilingRuntime)
		ctxt.DiagFunc = func(format string, args ...interface{}) {
			diag = true
			log.Printf(format, args...)
		}
		if *flags.SymABIs {
			ok = parser.ParseSymABIs(buf)
		} else {
			pList := new(obj.Plist)
			pList.Firstpc, ok = parser.Parse()
			// reports errors to parser.Errorf
			if ok {
				obj.Flushplist(ctxt, pList, nil, *flags.Importpath)
			}
		}
		if !ok {
			failedFile = f
			break
		}
	}
	if ok && !*flags.SymABIs {
		ctxt.NumberSyms()
		obj.WriteObjFile(ctxt, buf)
	}
	if !ok || diag {
		if failedFile != "" {
			log.Printf("assembly of %s failed", failedFile)
		} else {
			log.Print("assembly failed")
		}
		buf.Close()
		os.Remove(*flags.OutputFile)
		os.Exit(1)
	}
}
