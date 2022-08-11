// Copyright 2022 Bergur Ragnarsson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ioutil

import (
	"io"
	"io/ioutil"
	"os"

	"r2.is/gocompiler/vfs"
	vos "r2.is/gocompiler/vfs/os"
	"strings"
)

func NopCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}
func ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

func ReadFile(filename string) ([]byte, error) {
	filename2 := vfs.CleanPath(filename)
	if strings.HasPrefix(filename2, vfs.GorootSrc) {
		return vfs.ReadFile(filename2)
	}
	return ioutil.ReadFile(filename)
}
func ReadDir(dirname string) ([]os.FileInfo, error) {
	dirname2 := vfs.CleanPath(dirname)
	if strings.HasPrefix(dirname2, vfs.GorootSrc) {
		dirname = dirname2
		dlist, err := vfs.ReadDir(dirname)
		if err != nil {
			return nil, err
		}
		var ilist []os.FileInfo
		for _, v := range dlist {
			info, err := v.Info()
			if err != nil {
				return nil, err
			}
			ilist = append(ilist, info)
		}
		return ilist, nil
	} else {
		return ioutil.ReadDir(dirname)
	}
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

var Discard io.Writer = ioutil.Discard //devNull(0)

func TempFile(dir, pattern string) (f *vos.File, err error) {
	return vos.Wrap(ioutil.TempFile(dir, pattern))
}

func TempDir(dir, pattern string) (name string, err error) {
	return ioutil.TempDir(dir, pattern)
}
