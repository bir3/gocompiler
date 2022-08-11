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

        "r2.is/gocompiler/src/cmd/asm" //syncimport: "$pkg/src/$cmd/asm"
        "r2.is/gocompiler/src/cmd/cgo" //syncimport: "$pkg/src/$cmd/cgo"
        "r2.is/gocompiler/src/cmd/compile" //syncimport: "$pkg/src/$cmd/compile"
        "r2.is/gocompiler/src/cmd/gocmd" //syncimport: "$pkg/src/$cmd/gocmd"
        "r2.is/gocompiler/src/cmd/link" //syncimport: "$pkg/src/$cmd/link"
        "r2.is/gocompiler/vfs" //syncimport: "$pkg/vfs"
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
	default:
		fmt.Fprintf(os.Stderr, "ERROR: unknown GOCOMPILER_TOOL=%s\n", os.Getenv("GOCOMPILER_TOOL"))
		os.Exit(3)
	}
}

type Result struct {
	Stdout string
	Stderr string
}

func RunWithEnv(env []string, args ...string) (Result, error) {
	var result Result

	if len(args) < 2 {
		return result, errors.New("too few arguments")
	}
	if args[0] != "go" {
		return result, errors.New("only 'go' supported as first argument")
	}
	cmd := exec.Command(os.Args[0], args[1:]...)

	for _, s := range env {
		cmd.Env = append(cmd.Env, s)
	}
	cmd.Env = append(cmd.Env, "GOCOMPILER_TOOL=go")

	var out bytes.Buffer
	var outerr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &outerr

	err := cmd.Run()
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
