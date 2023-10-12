// Copyright 2023 Bergur Ragnarsson
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extra

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileMode = os.FileMode

func MkdirAllRace(dir string, perm FileMode) error {
	// safe for many processes to run concurrently
	dir = filepath.Clean(dir)
	if !filepath.IsAbs(dir) {
		return fmt.Errorf("not absolute path: %s", dir)
	}
	missing, err := missingFolders(dir, []string{})
	if err != nil {
		return err
	}
	for _, d2 := range missing {
		os.Mkdir(d2, perm) // ignore error as we may race
	}

	// at the end, we want a folder to exist
	// - no matter who created it:
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("failed to create folder %s - %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a folder %s", dir)
	}

	return nil
}

func missingFolders(dir string, missing []string) ([]string, error) {
	for {
		info, err := os.Stat(dir)
		if err == nil {
			if info.IsDir() {
				return missing, nil
			}
			return []string{}, fmt.Errorf("not a folder: %s", dir)
		}
		missing = append([]string{dir}, missing...) // prepend => reverse order
		d2 := filepath.Dir(dir)
		if d2 == dir {
			break
		}
		dir = d2
	}
	return []string{}, fmt.Errorf("program error at folder: %s", dir)
}
