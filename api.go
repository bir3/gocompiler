// Copyright 2022 Bergur Ragnarsson.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocompiler //syncpackage:

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

        "github.com/bir3/gocompiler/src/cmd/asm" //syncimport: "$pkg/src/$cmd/asm"
        "github.com/bir3/gocompiler/src/cmd/cgo" //syncimport: "$pkg/src/$cmd/cgo"
        "github.com/bir3/gocompiler/src/cmd/compile" //syncimport: "$pkg/src/$cmd/compile"
        "github.com/bir3/gocompiler/src/cmd/gocmd" //syncimport: "$pkg/src/$cmd/gocmd"
        "github.com/bir3/gocompiler/src/cmd/gofmt" //syncimport: "$pkg/src/$cmd/gofmt"
        "github.com/bir3/gocompiler/src/cmd/link" //syncimport: "$pkg/src/$cmd/link"

	/*
        "github.com/bir3/gocompiler/src/cmd/addr2line" //syncimport: "$pkg/src/$cmd/addr2line"
        "github.com/bir3/gocompiler/src/cmd/buildid" //syncimport: "$pkg/src/$cmd/buildid"
        "github.com/bir3/gocompiler/src/cmd/covdata" //syncimport: "$pkg/src/$cmd/covdata"
        "github.com/bir3/gocompiler/src/cmd/dist" //syncimport: "$pkg/src/$cmd/dist"
        "github.com/bir3/gocompiler/src/cmd/doc" //syncimport: "$pkg/src/$cmd/doc"
        "github.com/bir3/gocompiler/src/cmd/fix" //syncimport: "$pkg/src/$cmd/fix"
        "github.com/bir3/gocompiler/src/cmd/nm" //syncimport: "$pkg/src/$cmd/nm"
        "github.com/bir3/gocompiler/src/cmd/objdump" //syncimport: "$pkg/src/$cmd/objdump"
        "github.com/bir3/gocompiler/src/cmd/pack" //syncimport: "$pkg/src/$cmd/pack"
        "github.com/bir3/gocompiler/src/cmd/test2json" //syncimport: "$pkg/src/$cmd/test2json"
        "github.com/bir3/gocompiler/src/cmd/trace" //syncimport: "$pkg/src/$cmd/trace"
	*/
        "github.com/bir3/gocompiler/vfs" //syncimport: "$pkg/vfs"
)

func DebugShowEmbed() {
	vfs.DebugShowEmbed()
}

func IsRunToolchainRequest() bool {
	return os.Getenv("GOCOMPILER_TOOL") != ""
}

// adding extra executables : 44 MB -> 51 MB
func RunToolchain() {
	switch os.Getenv("GOCOMPILER_TOOL") {
	/*
		case "addr2line":
			addr2line.Main()
		case "buildid":
			buildid.Main()
		case "covdata":
			covdata.Main()
		case "dist":
			dist.Main()
		case "doc":
			doc.Main()
		case "fix":
			fix.Main()
		case "nm":
			nm.Main()
		case "objdump":
			objdump.Main()
		case "pack":
			pack.Main()
		case "test2json":
			test2json.Main()
		case "trace":
			trace.Main()
	*/
	case "go":
		gocmd.Main()
	case "compile":
		compile.Main()
	case "asm":
		asm.Main()
	case "link":
		link.Main()
	case "cgo":
		cgo.Main()
	case "gofmt":
		gofmt.Main()
	default:
		fmt.Fprintf(os.Stderr, "ERROR: unknown GOCOMPILER_TOOL=%s\n", os.Getenv("GOCOMPILER_TOOL"))
		os.Exit(3)
	}
}

type Result struct {
	Stdout string
	Stderr string
}

func Command(env []string, args ...string) (*exec.Cmd, error) {
	if vfs.SharedExe == "" {
		if vfs.SharedExeError == nil {
			panic("program error")
		}
		return nil, vfs.SharedExeError
	}
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}

	cmd := exec.Command(vfs.SharedExe, args[1:]...)

	cmd.Env = make([]string, len(env), len(env)+10)
	copy(cmd.Env, env)

	cmd.Env = append(cmd.Env, fmt.Sprintf("GOCOMPILER_TOOL=%s", args[0]))

	return cmd, nil
}

func RunWithEnv(env []string, args ...string) (Result, error) {
	var result Result

	cmd, err := Command(env, args...)
	if err != nil {
		return result, err
	}
	var out bytes.Buffer
	var outerr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &outerr

	err = cmd.Run()
	result.Stdout = out.String()
	result.Stderr = outerr.String()

	if err != nil {
		return result, err
	}
	return result, nil
}

func Run(args ...string) (Result, error) {
	return RunWithEnv(os.Environ(), args...)
}

func GoVersion() string {
	return "go1.20.6"
}
