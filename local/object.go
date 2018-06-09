package local

import (
	"github.com/usmanhalalit/gost"
	"os"
)

type Object struct {
	Path string
	Fs   *LocalFilesystem
}

func (f *Object) GetPath() string {
	return f.Path
}

func (f *Object) String() string {
	return f.GetPath()
}

func (f *Object) Exist() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *Object) Delete() error {
	return os.Remove(f.Path)
}

func (f *Object) Stat() (gost.FileInfo, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return gost.FileInfo{}, err
	}

	return gost.FileInfo{
		Size: fi.Size(),
		LastModified: fi.ModTime(),
	}, nil
}
