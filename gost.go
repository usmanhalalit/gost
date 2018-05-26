package gost

import "io"

type ReadWriteCloseSeeker interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker
}