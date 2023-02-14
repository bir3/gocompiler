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
	"strings"

        "github.com/bir3/gocompiler/src/cmd/asm" //syncimport: "$pkg/src/$cmd/asm"
        "github.com/bir3/gocompiler/src/cmd/cgo" //syncimport: "$pkg/src/$cmd/cgo"
        "github.com/bir3/gocompiler/src/cmd/compile" //syncimport: "$pkg/src/$cmd/compile"
        "github.com/bir3/gocompiler/src/cmd/gocmd" //syncimport: "$pkg/src/$cmd/gocmd"
        "github.com/bir3/gocompiler/src/cmd/gofmt" //syncimport: "$pkg/src/$cmd/gofmt"
        "github.com/bir3/gocompiler/src/cmd/link" //syncimport: "$pkg/src/$cmd/link"
        "github.com/bir3/gocompiler/vfs" //syncimport: "$pkg/vfs"
)

func DebugShowEmbed() {
	vfs.DebugShowEmbed()
}

func IsRunToolchainRequest() bool {
	return os.Getenv("GOCOMPILER_TOOL") != ""
}

func RunToolchain() {
	switch os.Getenv("GOCOMPILER_TOOL") {
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
	if !(args[0] == "go" || args[0] == "gofmt") {
		return nil, errors.New("only 'go' or 'gofmt' supported as first argument")
	}
	cmd := exec.Command(vfs.SharedExe, args[1:]...)

	cmd.Env = make([]string, len(env), len(env)+10)
	copy(cmd.Env, env)

	// disable cgo for now, does not work yet
	var cgoVar bool = false
	cmd.Env = append(cmd.Env, "GOCOMPILER_TOOL=go")
	for _, s := range cmd.Env {
		if strings.HasPrefix(s, "CGO_ENABLED=") {
			cgoVar = true
		}
	}
	if !cgoVar {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	}
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
	return "1.19.6"
}
