// Copyright 2022 Bergur Ragnarsson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import "io"

type Reader = io.Reader
type ReadCloser = io.ReadCloser
type ReadWriter = io.ReadWriter // go vet
type Writer = io.Writer
type Closer = io.Closer
type WriteCloser = io.WriteCloser
type ReadSeeker = io.ReadSeeker
type ByteWriter = io.ByteWriter
type WriterAt = io.WriterAt

type ReaderAt = io.ReaderAt
type SectionReader = io.SectionReader
type LimitedReader = io.LimitedReader

const (
	SeekStart   = io.SeekStart   // seek relative to the origin of the file
	SeekCurrent = io.SeekCurrent // seek relative to the current offset
	SeekEnd     = io.SeekEnd     // seek relative to the end
)

var EOF = io.EOF
var ErrUnexpectedEOF = io.ErrUnexpectedEOF
var ErrNoProgress = io.ErrNoProgress

func Copy(dst Writer, src Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func WriteString(w Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

func ReadFull(r Reader, buf []byte) (n int, err error) {
	return io.ReadFull(r, buf)
}

func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader {
	return io.NewSectionReader(r, off, n)
}

func CopyN(dst Writer, src Reader, n int64) (written int64, err error) {
	return io.CopyN(dst, src, n)
}

func MultiWriter(writers ...Writer) Writer {
	return io.MultiWriter(writers...)
}

func ReadAll(r Reader) ([]byte, error) {
	return io.ReadAll(r)
}

var Discard Writer = io.Discard

var LimitReader = io.LimitReader

func NopCloser(r Reader) ReadCloser { return io.NopCloser(r) }
