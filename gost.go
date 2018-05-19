package gost

import "io"

type ReadWriteCloser interface {
	io.Reader
	io.Writer
	io.Closer
}