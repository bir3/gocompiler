// Copyright 2022 Bergur Ragnarsson.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfs

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bir3/gocompiler/extra"
	"github.com/bir3/gocompiler/extra/extract_stdlib"
)

//go:embed goroot
var content embed.FS

var GOROOT = "/github.com/bir3/gocompiler/error-missing-init"

func init() {
	if os.Getenv("BIR3_GOCOMPILER_TOOL") == "" {
		return // not in Go toolchain mode
	}

	d, err := PrivateGOROOT()
	if err != nil {
		return // compiler will fail due to missing GOROOT
	}
	GOROOT = d
}

func SetupStdlib() error {
	d, err := PrivateGOROOT()
	if err != nil {
		return err
	}
	err = extractStdlib(d) // no-op if already done
	if err != nil {
		return err
	}
	return nil
}

func PrivateGOROOT() (string, error) {
	return configDir("bir3-gocompiler/stdlib-go1.22.0-0681") //syncvar:

}

func extractStdlib(d string) error {
	f, err := content.Open("goroot/stdlib-go1.22.0-0681.tar.zst") //syncvar:
	if err != nil {
		panic(fmt.Sprintf("gocompiler stdlib init failed - %s", err))
	}
	defer f.Close()

	err = extract_stdlib.ExtractStdlib(f, d)
	if err != nil {
		return fmt.Errorf("github.com/bir3/gocompiler: extract stdlib to %s failed with %v\n", d, err)
	}
	return nil
}

func configDir(folder string) (string, error) {
	var err error
	d := os.Getenv("BIR3_GOCOMPILER_GOROOT")
	if d == "" {
		d, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get config folder for stdlib - %w", err)
		}
	}
	d = filepath.Join(d, folder)
	if !filepath.IsAbs(d) {
		return "", fmt.Errorf("config folder %s is not absolute path for stdlib - %w", d, err)
	}
	if err == nil {
		err = extra.MkdirAllRace(d, 0755)
	}
	if err != nil {
		return "", fmt.Errorf("failed to create config folder %s for stdlib - %w", d, err)
	}
	readme := `
created by github.com/bir3/gocompiler

this folder can be set with env BIR3_GOCOMPILER_GOROOT
`
	os.WriteFile(filepath.Join(d, "README-bir3"), []byte(readme), 0666)
	return d, nil
}
