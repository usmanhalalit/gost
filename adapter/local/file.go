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
	path   string
	Fs     *LocalFilesystem
	reader io.Reader
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

func (f *File) Delete() error {
	return os.Remove(f.path)
}

func (f *File) Exist() bool {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *File) Stat() (adapter.FileInfo, error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return adapter.FileInfo{}, err
	}

	return adapter.FileInfo{
		Size: fi.Size(),
		LastModified: fi.ModTime(),
	}, nil
}

func (f *File) Directory() adapter.Directory {
	return &Directory{
		Path: filepath.Dir(f.path),
		Fs: f.Fs,
	}
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) Filesystem() adapter.Filesystem {
	return f.Fs
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		r, err := os.Open(f.path)
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

func (f *File) String() string {
	return f.GetPath()
}
