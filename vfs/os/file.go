// Copyright 2022 Bergur Ragnarsson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"r2.is/gocompiler/vfs"
	"strings"
	"time"
)

var Args []string = os.Args

type File struct {
	osfile  *os.File
	vfsfile fs.File
}

// special case for gocompiler/cmdlib/cgolib/out.go
func (f File) Getosfile() *os.File {
	return f.osfile
}

type FileInfo = os.FileInfo
type PathError = os.PathError
type LinkError = os.LinkError
type FileMode = os.FileMode

type Signal = os.Signal

type SyscallError = os.SyscallError

const (
	SEEK_SET int = os.SEEK_SET //= 0 // seek relative to the origin of the file
	SEEK_CUR int = os.SEEK_CUR // 1 // seek relative to the current offset
	SEEK_END int = os.SEEK_END
)

const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = os.O_RDONLY //syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = os.O_WRONLY // syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = os.O_RDWR   //syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = os.O_APPEND //syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = os.O_CREATE //syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = os.O_EXCL   //syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = os.O_SYNC   //syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = os.O_TRUNC  //syscall.O_TRUNC  // truncate regular writable file when opened.
)

const (
	ModeSymlink = os.ModeSymlink
	ModeSetgid  = os.ModeSetgid
	ModeType    = os.ModeType
)
const ModePerm FileMode = os.ModePerm

const DevNull = os.DevNull

var PathSeparator = os.PathSeparator

var ErrNotExist = os.ErrNotExist
var ErrClosed = os.ErrClosed

var Interrupt = os.Interrupt

/* if Stdin are not pure *os.File then Cmd.Wait fails to detect
them as such and will try to read Stdin, causing a process to hang
until stdin stream is closed (e.g. via control-D)
*/

var (
	Stdin  = os.Stdin
	Stdout = os.Stdout
	Stderr = os.Stderr
)

//------------------------------------

func wrap(f *os.File, err error) (*File, error) {
	return &File{f, nil}, err
}

// for ioutil:
func Wrap(f *os.File, err error) (*File, error) {
	return wrap(f, err)
}

//----------------------------------
func errIfVfsfile(f File, where string) {
	if f.vfsfile != nil {
		panic("vfsfile not supported here yet: " + where)
	}
}

func Chdir(dir string) error {
	dir2 := vfs.CleanPath(dir)
	if strings.HasPrefix(dir2, vfs.GorootSrc) {
		// vfsfs Chdir
		return vfs.Chdir(dir2)
	} else {
		return os.Chdir(dir)
	}
}

func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	var f File
	var err error
	name2 := vfs.CleanPath(name)
	if strings.HasPrefix(name2, vfs.GOROOT) {
		name = name2
		// our internal goroot
		// BUG: we ignore flag and perm (should error if can't satisfy)
		f.vfsfile, err = vfs.Open(name)
	} else {
		f.osfile, err = os.OpenFile(name, flag, perm)
	}

	if vfs.LogFile != nil {
		var zstatus string = "- "
		if strings.HasPrefix(name2, vfs.GOROOT) {
			zstatus = "vfs    "
			if err != nil {
				zstatus = "vfs-err"
			}
		} else {
			zstatus = "os     "
			if err != nil {
				zstatus = "os-err "
			}
		}
		vfs.Log2(err, zstatus+" Open "+name)
	}
	return &f, err
}

func Remove(name string) error {
	return os.Remove(name)
}

func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func Create(name string) (*File, error) {
	return wrap(os.Create(name))
}

func Stat(name string) (FileInfo, error) {
	name2 := vfs.CleanPath(name)
	if strings.HasPrefix(name2, vfs.GOROOT) {
		return vfs.Stat(name2)
	} else {
		fi, err := os.Stat(name)
		if vfs.LogFile != nil {
			vfs.Log2(err, "Stat "+name)
		}
		return fi, err
	}
}

func Lstat(name string) (FileInfo, error) {
	name2 := vfs.CleanPath(name)
	if strings.HasPrefix(name2, vfs.GOROOT) {
		return vfs.Stat(name2)
	}
	return os.Lstat(name)
}

func Readlink(name string) (string, error) {
	name2 := vfs.CleanPath(name)
	if strings.HasPrefix(name2, vfs.GOROOT) {
		return "", &os.PathError{Op: "readlink", Path: name, Err: errors.New("gocompiler/vfs path")}
	}
	return os.Readlink(name)
}

func Symlink(oldname, newname string) error {
	oldname2 := vfs.CleanPath(oldname)
	newname2 := vfs.CleanPath(newname)
	if strings.HasPrefix(oldname2, vfs.GOROOT) {
		return &LinkError{"symlink", oldname, newname, errors.New("gocompiler/vfs path")}
	}
	if strings.HasPrefix(newname2, vfs.GOROOT) {
		return &LinkError{"symlink", oldname, newname, errors.New("gocompiler/vfs path")}
	}
	return os.Symlink(oldname, newname)
}

func IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func IsExist(err error) bool {
	return os.IsExist(err)
}

func IsPermission(err error) bool {
	return os.IsPermission(err)
}

func Mkdir(name string, perm FileMode) error {
	return os.Mkdir(name, perm)
}

func MkdirAll(path string, perm FileMode) error {
	return os.MkdirAll(path, perm)
}

func Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func Chmod(name string, mode FileMode) error {
	return os.Chmod(name, mode)
}

func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func Exit(code int) {
	os.Exit(code)
}

func Getenv(key string) string {
	return os.Getenv(key)
}

func Environ() []string {
	return os.Environ()
}

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func Expand(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}

func ExpandEnv(s string) string {
	return os.ExpandEnv(s)
}

func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func Getwd() (dir string, err error) {
	return os.Getwd()
}

func Getuid() int {
	return os.Getuid()
}
func Getpagesize() int {
	return os.Getpagesize()
}
func TempDir() string {
	return os.TempDir()
}

func Executable() (string, error) {
	return os.Executable()
}

func SameFile(fi1, fi2 FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

func UserHomeDir() (string, error) {
	return os.UserHomeDir()
}
func UserConfigDir() (string, error) {
	return os.UserConfigDir()
}
func UserCacheDir() (string, error) {
	return os.UserCacheDir()
}

///-------------
func (f File) Read(b []byte) (n int, err error) {
	if f.osfile != nil {
		return f.osfile.Read(b)
	} else {
		return f.vfsfile.Read(b)
	}
}
func (f File) Seek(offset int64, whence int) (ret int64, err error) {
	errIfVfsfile(f, "Seek")
	return f.osfile.Seek(offset, whence)
}
func (f File) Close() error {
	if f.osfile != nil {
		return f.osfile.Close()
	} else {
		return f.vfsfile.Close()
	}
}
func (f File) Name() string {
	if f.osfile != nil {
		return f.osfile.Name()
	} else {
		fi, err := f.vfsfile.Stat()
		if err != nil {
			panic("vfs stat failed")
		}
		return fi.Name()
	}
}
func (f File) ReadAt(b []byte, off int64) (n int, err error) {
	errIfVfsfile(f, "ReadAt")
	return f.osfile.ReadAt(b, off)
}
func (f File) Stat() (FileInfo, error) {
	if f.osfile != nil {
		return f.osfile.Stat()
	} else {
		return f.vfsfile.Stat()
	}
}
func (f File) Readdirnames(n int) (names []string, err error) {
	errIfVfsfile(f, "Readdirnames")
	return f.osfile.Readdirnames(n)
}

func (f File) Fd() uintptr               { return f.osfile.Fd() }
func (f File) Truncate(size int64) error { return f.osfile.Truncate(size) }
func (f File) Write(b []byte) (n int, err error) {
	return f.osfile.Write(b)
}
func (f File) WriteAt(b []byte, off int64) (n int, err error) { return f.osfile.WriteAt(b, off) }
func (f File) Sync() error                                    { return f.osfile.Sync() }
func (f File) WriteString(s string) (n int, err error)        { return f.osfile.WriteString(s) }

func (f File) ReadDir(n int) ([]fs.DirEntry, error) {
	// bug: only osfile supported
	return f.osfile.ReadDir(n)
}

func CreateTemp(dir, pattern string) (*File, error) { return wrap(os.CreateTemp(dir, pattern)) }

func WriteFile(name string, data []byte, perm FileMode) error { return os.WriteFile(name, data, perm) }

func ReadDir(name string) ([]fs.DirEntry, error) {
	name2 := vfs.CleanPath(name)
	if strings.HasPrefix(name2, vfs.GOROOT) {
		return nil, errors.New("gocompiler vfs path")
	}

	return os.ReadDir(name)
}

func MkdirTemp(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func ReadFile(name string) ([]byte, error) {
	f, err := Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}

func IsPathSeparator(c uint8) bool {
	return os.IsPathSeparator(c)
}
