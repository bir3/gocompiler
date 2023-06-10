package extract_stdlib

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/bir3/gocompiler/src/cmd/gocmd/compress/zstd"
	"github.com/bir3/gocompiler/src/cmd/gocmd/lockedfile"
)

/*
 size go stdlib tarfile:
   37 MB full-notest.tar
   			# no *_test.go and testdata folders (else 56 MB)
			# no src/cmd
			# version: go1.20

-> compressed with zstd v1.5.5 :
	7.9 MB full-notest-3.zst    x4.7 compression   # -3 --ultra
 	5.7 MB full-notest-21.zst   x6.5 compression   # -21 --ultra

-> decompression: 76 milliseconds on macbook air m1
*/

func ExtractStdlib(file fs.File, dir string) error {

	donefile := path.Join(dir, "done")
	if _, err := os.Stat(donefile); err == nil {
		return nil // already extracted
	}

	lockfile := path.Join(dir, "done.lock")
	lf, err := lockedfile.Create(lockfile)
	defer lf.Close()
	if err != nil {
		return fmt.Errorf("Error creating lockfile: %w", err)
	}

	zstdReader, err := zstd.NewReader(file)
	if err != nil {
		return fmt.Errorf("Error creating zstd reader: %w", err)
	}
	defer zstdReader.Close()

	tarReader := tar.NewReader(zstdReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Errorf("Error reading tar header - %w", err)
		}

		target := path.Join(dir, header.Name)
		switch header.Typeflag {
		case tar.TypeReg: // regular file
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("Error creating directory %s - %w", target, err)
			}
			continue
		default:
			return fmt.Errorf("unknown tar item type %d", header.Typeflag)
		}

		// TODO: if python script creates "tar.TypeDir" entry, we can avoid
		//		next stat/mkdir

		// Create any necessary parent directories for the file
		parent := filepath.Dir(target)
		if _, err := os.Stat(parent); os.IsNotExist(err) {
			if err := os.MkdirAll(parent, 0755); err != nil {
				return fmt.Errorf("Error creating directory %s - %w", parent, err)
			}
		}

		// Extract the file
		outFile, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("Error creating file %s - %w", target, err)
		}

		if _, err := io.Copy(outFile, tarReader); err != nil {
			return fmt.Errorf("Error extracting file %s - %w", target, err)
		}

		err = outFile.Close()
		if err != nil {
			return fmt.Errorf("Error closing file %s - %w", target, err)
		}
		// note: we do not replicate chmod (e.g. executable flag)
	}
	f3, err := os.Create(donefile)
	if err != nil {
		return fmt.Errorf("Error creating donefile: %s - %w", donefile, err)
	}
	err = f3.Close()
	if err != nil {
		return fmt.Errorf("Error closing donefile: %s - %w", donefile, err)
	}
	return nil
}
