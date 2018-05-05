package local

import (
	"github.com/usmanhalalit/gost/adapter"
	"os"
)

type Object struct {
	Path string
	Fs   *LocalFilesystem
}

func (f *Object) Filesystem() adapter.Filesystem {
	return f.Fs
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

func (f *Object) Stat() (adapter.FileInfo, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return adapter.FileInfo{}, err
	}

	return adapter.FileInfo{
		Size: fi.Size(),
		LastModified: fi.ModTime(),
	}, nil
}
