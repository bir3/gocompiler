// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd
// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd

package base

import (
	       "r2.is/gocompiler/vfs/io"
	       "r2.is/gocompiler/vfs/os"
)

func MapFile(f *os.File, offset, length int64) (string, error) {
	buf := make([]byte, length)
	_, err := io.ReadFull(io.NewSectionReader(f, offset, length), buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
