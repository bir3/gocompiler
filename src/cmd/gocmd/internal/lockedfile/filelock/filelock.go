package filelock

import "github.com/bir3/gocompiler/src/cmd/gocmd/internal/lockedfile/internal/filelock"

// export 1
type File = filelock.File

func Lock(f File) error {
	return filelock.Lock(f)
}

func RLock(f File) error {
	return filelock.RLock(f)
}

func Unlock(f File) error {
	return filelock.Unlock(f)
}
