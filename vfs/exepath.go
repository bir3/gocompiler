// Copyright 2022 Bergur Ragnarsson.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vfs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func exepath(selfpath string) (string, error) {
	if len(selfpath) == 0 {
		return "", fmt.Errorf("empty string in os.Args[0]")
	}

	if path.IsAbs(selfpath) {
		return selfpath, nil
	}

	foundpath, err := exec.LookPath(selfpath)
	if err != nil {
		return "", err
	}
	if path.IsAbs(foundpath) {
		return foundpath, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(wd, foundpath), nil
}
