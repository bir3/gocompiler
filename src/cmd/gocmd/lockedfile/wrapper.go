package lockedfile

import (
	"io"
	"io/fs"

	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/lockedfile"
)

type File struct {
	lf *lockedfile.File
}

func mapfile(f *lockedfile.File) *File {
	if f != nil {
		return &File{f}
	}
	return nil
}

// Open is like os.Open, but returns a read-locked file.
func Open(name string) (*File, error) {
	f, err := lockedfile.Open(name)
	return mapfile(f), err
}

// Create is like os.Create, but returns a write-locked file.
func Create(name string) (*File, error) {
	f, err := lockedfile.Create(name)
	return mapfile(f), err
}

// Close unlocks and closes the underlying file.
//
// Close may be called multiple times; all calls after the first will return a
// non-nil error.
func (f *File) Close() error {
	return f.lf.Close()
}

// Write opens the named file (creating it with the given permissions if needed),
// then write-locks it and overwrites it with the given content.
func Write(name string, content io.Reader, perm fs.FileMode) (err error) {
	return lockedfile.Write(name, content, perm)
}

// Transform invokes t with the result of reading the named file, with its lock
// still held.
//
// If t returns a nil error, Transform then writes the returned contents back to
// the file, making a best effort to preserve existing contents on error.
//
// t must not modify the slice passed to it.
func Transform(name string, t func([]byte) ([]byte, error)) (err error) {
	return lockedfile.Transform(name, t)
}
