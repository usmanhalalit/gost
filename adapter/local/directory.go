package local

import (
	"github.com/usmanhalalit/gost/adapter"
	"io/ioutil"
	"strings"
)

type Directory struct {
	Path string
	Fs *LocalFilesystem
}

func (d *Directory) Filesystem() adapter.Filesystem {
	return d.Fs
}

func (d *Directory) File(path string) adapter.File {
	return &File{
		path:   d.Path + "/" + path,
		Fs:     d.Fs,
		reader: nil,
	}
}

func (d *Directory) GetPath() string {
	return d.Path
}

func (d *Directory) Directory(path string) adapter.Directory {
	path = d.Path + "/" + path
	path = strings.TrimRight(path, "/")
	return &Directory{
		Path: path,
		Fs: d.Fs,
	}
}

func (d *Directory) Files() ([]adapter.File, error) {
	files, err := ioutil.ReadDir(d.Path)
	if err != nil { return nil, err }
	var localFiles []adapter.File
	for i := range files {
		s3file := File{
			path:   d.Path + "/" + files[i].Name(),
			Fs:     d.Fs,
			reader: nil,
		}
		localFiles = append(localFiles, adapter.File(&s3file))
	}
	return localFiles, nil
}