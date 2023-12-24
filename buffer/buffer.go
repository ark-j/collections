package buffer

import (
	"io"
)

type Buffer struct {
	buf []byte // buffer of bytes
	off int    // buffer offset
}

func New(buf []byte) *Buffer { return &Buffer{buf: buf} }

func (b *Buffer) Write(p []byte) (n int, err error) {}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {}

func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {}

func (b *Buffer) Len() int {}
