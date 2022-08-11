// Copyright 2022 Bergur Ragnarsson.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfs

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"strings"
)

//go:embed goroot
var content embed.FS

/*
	this file should only contain general virtual filesystem things
	; the logic merging the normal "io" fs and this in-memory fs should live elsewhere
*/

func DebugShowEmbed() {
	fs.WalkDir(content, "goroot", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(path)
		return nil
	})
}

// most paths are expected to be clean, so favor that case
func CleanPath(s string) string {
	state := 0
	for _, c := range s {
		switch state {
		case 0:
			if c == '/' {
				state = 1
			}
		case 1:
			if c == '/' || c == '.' {
				return filepath.Clean(s)
			} else {
				state = 0
			}
		}
	}
	return s
}

func Open(name string) (fs.File, error) {
	// TODO: support cwd

	if strings.HasPrefix(name, GOROOT) {
		name = name[gorootPrefixLen:]
	}
	if e, err := content.Open(name); err == nil {
		return e, nil
	}
	return nil, fmt.Errorf("s2: file %s not found", name)
}

func Stat(name string) (fs.FileInfo, error) {

	if strings.HasPrefix(name, GorootTool) {
		// pretend tool file exist
		b := filepath.Base(name)
		if b == "asm" || b == "compile" || b == "link" || b == "cgo" {
			name = filepath.Join(GorootSrc, "flag/flag.go")
		}
	}
	if strings.HasPrefix(name, GOROOT) {
		name = name[gorootPrefixLen:]
	}

	if f, err := content.Open(name); err == nil {
		fi, err := f.Stat()
		f.Close()
		if err == nil {
			return fi, nil
		} else {
			return fi, err
		}
	} else {
		return nil, err
	}
}

func ReadDir(name string) ([]fs.DirEntry, error) {

	if strings.HasPrefix(name, GOROOT) {
		name = name[gorootPrefixLen:]
	}

	return content.ReadDir(name)
}

func ReadFile(filename string) ([]byte, error) {
	if strings.HasPrefix(filename, GOROOT) {
		filename = filename[gorootPrefixLen:]
	}
	return content.ReadFile(filename)
}
