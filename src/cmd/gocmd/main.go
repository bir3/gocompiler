// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ./mkalldocs.sh

package gocmd

import (
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/workcmd"
	"context"
	"github.com/bir3/gocompiler/src/cmd/gocmd/flag"
	"fmt"
	"github.com/bir3/gocompiler/src/internal/buildcfg"
	"log"
	"os"
	"path/filepath"
	"github.com/bir3/gocompiler/vfs"
	"strings"

	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/base"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/bug"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/cfg"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/clean"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/doc"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/envcmd"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/fix"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/fmtcmd"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/generate"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/get"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/help"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/list"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/modcmd"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/modfetch"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/modget"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/modload"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/run"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/test"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/tool"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/trace"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/version"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/vet"
	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/work"
)

func init() {
	base.Go.Commands = []*base.Command{
		bug.CmdBug,
		work.CmdBuild,
		clean.CmdClean,
		doc.CmdDoc,
		envcmd.CmdEnv,
		fix.CmdFix,
		fmtcmd.CmdFmt,
		generate.CmdGenerate,
		modget.CmdGet,
		work.CmdInstall,
		list.CmdList,
		modcmd.CmdMod,
		workcmd.CmdWork,
		run.CmdRun,
		test.CmdTest,
		tool.CmdTool,
		version.CmdVersion,
		vet.CmdVet,

		help.HelpBuildConstraint,
		help.HelpBuildmode,
		help.HelpC,
		help.HelpCache,
		help.HelpEnvironment,
		help.HelpFileType,
		modload.HelpGoMod,
		help.HelpGopath,
		get.HelpGopathGet,
		modfetch.HelpGoproxy,
		help.HelpImportPath,
		modload.HelpModules,
		modget.HelpModuleGet,
		modfetch.HelpModuleAuth,
		help.HelpPackages,
		modfetch.HelpPrivate,
		test.HelpTestflag,
		test.HelpTestfunc,
		modget.HelpVCS,
	}
}

func Main() {
	_ = go11tag
	flag.Usage = base.Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		base.Usage()
	}

	if args[0] == "get" || args[0] == "help" {
		if !modload.WillBeEnabled() {
			// Replace module-aware get with GOPATH get if appropriate.
			*modget.CmdGet = *get.CmdGet
		}
	}

	cfg.CmdName = args[0] // for error messages
	if args[0] == "help" {
		help.Help(os.Stdout, args[1:])
		return
	}

	// Diagnose common mistake: GOPATH==GOROOT.
	// This setting is equivalent to not setting GOPATH at all,
	// which is not what most people want when they do it.
	if gopath := cfg.BuildContext.GOPATH; filepath.Clean(gopath) == filepath.Clean(vfs.GOROOT) {
		fmt.Fprintf(os.Stderr, "warning: GOPATH set to GOROOT (%s) has no effect\n", gopath)
	} else {
		for _, p := range filepath.SplitList(gopath) {
			// Some GOPATHs have empty directory elements - ignore them.
			// See issue 21928 for details.
			if p == "" {
				continue
			}
			// Note: using HasPrefix instead of Contains because a ~ can appear
			// in the middle of directory elements, such as /tmp/git-1.8.2~rc3
			// or C:\PROGRA~1. Only ~ as a path prefix has meaning to the shell.
			if strings.HasPrefix(p, "~") {
				fmt.Fprintf(os.Stderr, "go: GOPATH entry cannot start with shell metacharacter '~': %q\n", p)
				os.Exit(2)
			}
			if !filepath.IsAbs(p) {
				if cfg.Getenv("GOPATH") == "" {
					// We inferred $GOPATH from $HOME and did a bad job at it.
					// Instead of dying, uninfer it.
					cfg.BuildContext.GOPATH = ""
				} else {
					fmt.Fprintf(os.Stderr, "go: GOPATH entry is relative; must be absolute path: %q.\nFor more details see: 'go help gopath'\n", p)
					os.Exit(2)
				}
			}
		}
	}

	if cfg.GOROOT == "" {
		fmt.Fprintf(os.Stderr, "go: cannot find GOROOT directory: 'go' binary is trimmed and GOROOT is not set\n")
		os.Exit(2)
	}
	if fi, err := os.Stat(cfg.GOROOT); err != nil || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "go: cannot find GOROOT directory: %v\n", cfg.GOROOT)
		os.Exit(2)
	}

BigCmdLoop:
	for bigCmd := base.Go; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					help.PrintUsage(os.Stderr, bigCmd)
					base.SetExitStatus(2)
					base.Exit()
				}
				if args[0] == "help" {
					// Accept 'go mod help' and 'go mod help foo' for 'go help mod' and 'go help mod foo'.
					help.Help(os.Stdout, append(strings.Split(cfg.CmdName, " "), args[1:]...))
					return
				}
				cfg.CmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			invoke(cmd, args)
			base.Exit()
			return
		}
		helpArg := ""
		if i := strings.LastIndex(cfg.CmdName, " "); i >= 0 {
			helpArg = " " + cfg.CmdName[:i]
		}
		fmt.Fprintf(os.Stderr, "go %s: unknown command\nRun 'go help%s' for usage.\n", cfg.CmdName, helpArg)
		base.SetExitStatus(2)
		base.Exit()
	}
}

func invoke(cmd *base.Command, args []string) {
	// 'go env' handles checking the build config
	if cmd != envcmd.CmdEnv {
		buildcfg.Check()
		if cfg.ExperimentErr != nil {
			base.Fatalf("go: %v", cfg.ExperimentErr)
		}
	}

	// Set environment (GOOS, GOARCH, etc) explicitly.
	// In theory all the commands we invoke should have
	// the same default computation of these as we do,
	// but in practice there might be skew
	// This makes sure we all agree.
	cfg.OrigEnv = os.Environ()
	cfg.CmdEnv = envcmd.MkEnv()
	for _, env := range cfg.CmdEnv {
		if os.Getenv(env.Name) != env.Value {
			os.Setenv(env.Name, env.Value)
		}
	}

	cmd.Flag.Usage = func() { cmd.Usage() }
	if cmd.CustomFlags {
		args = args[1:]
	} else {
		base.SetFromGOFLAGS(&cmd.Flag)
		cmd.Flag.Parse(args[1:])
		args = cmd.Flag.Args()
	}
	ctx := maybeStartTrace(context.Background())
	ctx, span := trace.StartSpan(ctx, fmt.Sprint("Running ", cmd.Name(), " command"))
	cmd.Run(ctx, cmd, args)
	span.Done()
}

func init() {
	base.Usage = mainUsage
}

func mainUsage() {
	help.PrintUsage(os.Stderr, base.Go)
	os.Exit(2)
}

func maybeStartTrace(pctx context.Context) context.Context {
	if cfg.DebugTrace == "" {
		return pctx
	}

	ctx, close, err := trace.Start(pctx, cfg.DebugTrace)
	if err != nil {
		base.Fatalf("failed to start trace: %v", err)
	}
	base.AtExit(func() {
		if err := close(); err != nil {
			base.Fatalf("failed to stop trace: %v", err)
		}
	})

	return ctx
}
