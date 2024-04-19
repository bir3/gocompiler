// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go work init

package workcmd

import (
	"context"
	"path/filepath"

	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/base"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/fsys"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/gover"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/modload"

	"github.com/bir3/gocompiler/src/xvendor/golang.org/x/mod/modfile"
)

var cmdInit = &base.Command{
	UsageLine:	"go work init [moddirs]",
	Short:		"initialize workspace file",
	Long: `Init initializes and writes a new go.work file in the
current directory, in effect creating a new workspace at the current
directory.

go work init optionally accepts paths to the workspace modules as
arguments. If the argument is omitted, an empty workspace with no
modules will be created.

Each argument path is added to a use directive in the go.work file. The
current go version will also be listed in the go.work file.

See the workspaces reference at https://go.dev/ref/mod#workspaces
for more information.
`,
	Run:	runInit,
}

func init() {
	base.AddChdirFlag(&cmdInit.Flag)
	base.AddModCommonFlags(&cmdInit.Flag)
}

func runInit(ctx context.Context, cmd *base.Command, args []string) {
	modload.InitWorkfile()

	modload.ForceUseModules = true

	gowork := modload.WorkFilePath()
	if gowork == "" {
		gowork = filepath.Join(base.Cwd(), "go.work")
	}

	if _, err := fsys.Stat(gowork); err == nil {
		base.Fatalf("go: %s already exists", gowork)
	}

	goV := gover.Local()	// Use current Go version by default
	wf := new(modfile.WorkFile)
	wf.Syntax = new(modfile.FileSyntax)
	wf.AddGoStmt(goV)
	workUse(ctx, gowork, wf, args)
	modload.WriteWorkFile(gowork, wf)
}
