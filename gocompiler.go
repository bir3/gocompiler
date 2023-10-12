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
	"path/filepath"

        "github.com/bir3/gocompiler/extra" //syncimport: "$pkg/extra"
        "github.com/bir3/gocompiler/src/cmd/asm" //syncimport: "$pkg/src/$cmd/asm"
        "github.com/bir3/gocompiler/src/cmd/cgo" //syncimport: "$pkg/src/$cmd/cgo"
        "github.com/bir3/gocompiler/src/cmd/compile" //syncimport: "$pkg/src/$cmd/compile"
        "github.com/bir3/gocompiler/src/cmd/gocmd" //syncimport: "$pkg/src/$cmd/gocmd"
        "github.com/bir3/gocompiler/src/cmd/gofmt" //syncimport: "$pkg/src/$cmd/gofmt"
        "github.com/bir3/gocompiler/src/cmd/link" //syncimport: "$pkg/src/$cmd/link"
        "github.com/bir3/gocompiler/vfs" //syncimport: "$pkg/vfs"
)

type Info struct {
	GoVersion string
	CacheDir  string
	GOROOT    string
}

func GetInfo() (Info, error) {
	info := Info{}
	info.GoVersion = GoVersion()
	d, err := cacheDir()
	if err != nil {
		return Info{}, err
	}
	info.CacheDir = d
	info.GOROOT, err = vfs.PrivateGOROOT()
	if err != nil {
		return Info{}, err
	}
	return info, nil
}

func IsRunToolchainRequest() bool {
	return os.Getenv("BIR3_GOCOMPILER_TOOL") != ""
}

// adding extra executables : 44 MB -> 51 MB
func RunToolchain() {
	switch os.Getenv("BIR3_GOCOMPILER_TOOL") {
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
	case "debug-info":
		info, err := GetInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}
		fmt.Printf("Embedded Go compiler github.com/bir3/gocompiler\n")
		fmt.Printf("go-version    : %s\n", info.GoVersion)
		fmt.Printf("cache-dir     : %s\n", info.CacheDir)
		fmt.Printf("GOROOT/stdlib : %s\n", info.GOROOT)
	default:
		fmt.Fprintf(os.Stderr, "ERROR: unknown BIR3_GOCOMPILER_TOOL=%s\n", os.Getenv("BIR3_GOCOMPILER_TOOL"))
		os.Exit(3)
	}
}

type Result struct {
	Stdout string
	Stderr string
}

func cacheDir() (string, error) {
	d, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	d = filepath.Join(d, "bir3-gocompiler")

	if !filepath.IsAbs(d) {
		return "", fmt.Errorf("not absolute path: %s", d)
	}

	info, err := os.Stat(d)
	if err == nil {
		if info.IsDir() {
			return d, nil
		}
		return "", fmt.Errorf("not a folder: %s", d)
	}
	// this runs only once
	err = extra.MkdirAllRace(d, 0777)
	if err != nil {
		return "", nil
	}
	readme := `
cache for github.com/bir3/gocompiler
= Go compiler as a package
= private Go build cache to avoid interfering with the standard Go toolchain build cache
`
	os.WriteFile(filepath.Join(d, "README-bir3"), []byte(readme), 0666)
	return d, nil
}

func Command(env []string, args ...string) (*exec.Cmd, error) {
	SharedExe, err := os.Executable()

	if err != nil {
		return nil, err
	}
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}
	err = vfs.SetupStdlib() // no-op if already done
	if err != nil {
		return nil, err
	}
	privateCacheDir, err := cacheDir()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(SharedExe, args[1:]...)

	cmd.Env = make([]string, len(env), len(env)+10)
	copy(cmd.Env, env)

	cmd.Env = append(cmd.Env, fmt.Sprintf("BIR3_GOCOMPILER_TOOL=%s", args[0]))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOCACHE=%s", privateCacheDir))
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
	return "go1.20.10"
}
