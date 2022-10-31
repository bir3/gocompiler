// Copyright 2022 Bergur Ragnarsson.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const GOROOT = "/_/github.com/bir3/gocompiler/vfs/goroot"
const gorootPrefixLen = len("/_/github.com/bir3/gocompiler/vfs/")

var GorootSrc string
var GorootTool string
var SharedExe string
var SharedExeError error

func init() {
	// ; cmd/go/internal/cfg/cfg.go

	GorootSrc = filepath.Join(GOROOT, "src") + string(os.PathSeparator)
	GorootTool = filepath.Join(GOROOT, "pkg", "tool") + string(os.PathSeparator)

	if !filepath.IsAbs(os.Args[0]) {
		wd, err := os.Getwd()
		if err == nil {
			SharedExe = filepath.Join(wd, os.Args[0])
		} else {
			SharedExeError = fmt.Errorf("os.Getwd() failed : %w", err)
		}
	} else {
		SharedExe = os.Args[0]
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

var LogFile *os.File

func StackTrace() string {
	b := make([]byte, 9900) // adjust buffer size to be larger than expected stack
	n := runtime.Stack(b, false)
	s := string(b[:n])
	return s
}

func Log(msg string) {
	if LogFile != nil {
		LogFile.WriteString(msg + "\n")
	}
}

func Log2(err error, msg string) {
	errstr := "ok  "
	if err != nil {
		errstr = "err "
	}
	Log(errstr + msg)
	if err != nil {
		Log(StackTrace())
	}
}

var vfsCwd string

func Chdir(dir string) error {
	if !strings.HasPrefix(dir, GorootSrc) {
		panic("vfs Chdir, only supported for GorootSrc, got dir=" + dir)
	}
	vfsCwd = dir
	if LogFile != nil {
		Log("vfs: vfsCwd set to " + vfsCwd + " exe=" + os.Args[0])
	}
	return nil
}

func zpath(path string) string {
	if vfsCwd == "" {
		return path // default
	}
	if strings.HasPrefix(path, "/") {
		return path // absolute path
	}
	return filepath.Join(vfsCwd, path)
}
