// Copyright 2022 Bergur Ragnarsson.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfs

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"strings"

	"github.com/bir3/gocompiler/src/cmd/gocmd/extract_stdlib"
)

//go:embed goroot
var content embed.FS

/*
	"vfs2" = historical name, not virtual file system
*/

var GOROOT = "/github.com/bir3/gocompiler/missing-init"

const gorootPrefixLen = len("/_/github.com/bir3/gocompiler/vfs/")

// var GorootSrc string
var GorootTool string

var SharedExe string // used by so we can run as "compile", "asm", etc.
var SharedExeError error

func init() {
	// ; cmd/go/internal/cfg/cfg.go

	/*
		user-config-dir/gocompiler/src-xxxx/done

		override with env GOCOMPILER_DIR/src-xxxx/done
	*/
	//var err error

	d, err := configDir("gocompiler/stdlib-go1.21rc2-77dd") //syncvar:
	if err != nil {
		return // compiler will fail due to missing GOROOT
	}
	f, err := content.Open("goroot/stdlib-go1.21rc2-77dd.tar.zst") //syncvar:
	if err != nil {
		panic(fmt.Sprintf("gocompiler stdlib init failed - %s", err))
	}
	defer f.Close()
	GOROOT = d

	err = extract_stdlib.ExtractStdlib(f, d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: gocompiler: extract stdlib to %s failed with %v\n", d, err)
		os.Exit(3)
	}
	//GorootSrc = filepath.Join(GOROOT, "src") + string(os.PathSeparator)
	GorootTool = filepath.Join(GOROOT, "pkg", "tool") + string(os.PathSeparator)

	SharedExe, err = os.Executable()
	if err != nil {
		SharedExeError = err
	}
	switch os.Getenv("GOCOMPILERLIB_LOG") {
	case "1":
		basename := os.Args[0]
		basename = basename[strings.LastIndex(basename, `/`)+1:]
		basename = basename[strings.LastIndex(basename, `\`)+1:]
		basename = strings.TrimSuffix(basename, ".exe")

		filename := fmt.Sprintf("/tmp/gocompilerlib-%s-%d", basename, os.Getpid())
		LogFile, _ = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666) //Create("/tmp/r2log")

	}
}

func configDir(folder string) (string, error) {
	var err error
	d := os.Getenv("GOCOMPILER_DIR")
	if d == "" {
		d, err = os.UserCacheDir()
	}
	d = path.Join(d, folder)
	if err == nil {
		err = os.MkdirAll(d, 0755)
	}
	if err != nil {
		return "", fmt.Errorf("failed to create folder for stdlib - %w", err)
	}
	return d, nil
}

var LogFile *os.File

func DebugShowEmbed() {
	fs.WalkDir(content, "goroot", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(path)
		return nil
	})
}
