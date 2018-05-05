package local

import (
	"errors"
	"fmt"
	"github.com/usmanhalalit/gost/adapter"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Object
	reader io.ReadCloser
}

func (f *File) ReadString() (string, error) {
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func (f *File) WriteString(s string) error {
	b := []byte(s)
	n, err := f.Write(b)
	if n != len(b) {
		return errors.New(fmt.Sprintf("Wrote %v bytes from given %v bytes", n, len(b)))
	}
	return err
}

func (f *File) Directory() adapter.Directory {
	return &Directory{
		Object: Object{
			Path: filepath.Dir(f.Path),
			Fs:   f.Fs,
		},
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		r, err := os.Open(f.Path)
		if err != nil {
			return 0, err
		}
		f.reader = r
	}

	return f.reader.Read(p)
}

func (f *File) Write(p []byte) (n int, err error) {
	file, err := os.Create(f.GetPath())
	n, err = file.Write(p)
	return n, err
}

func (f *File) Close() error {
	return f.reader.Close()
}
