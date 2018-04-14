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
		file := files[i]
		if file.IsDir() {
			continue
		}
		localFile := File{
			path:   d.Path + "/" + file.Name(),
			Fs:     d.Fs,
			reader: nil,
		}
		localFiles = append(localFiles, adapter.File(&localFile))
	}
	return localFiles, nil
}

func (d *Directory) Directories() ([]adapter.Directory, error) {
	files, err := ioutil.ReadDir(d.Path)
	if err != nil { return nil, err }
	var localDirs []adapter.Directory
	for i := range files {
		dir := files[i]
		if ! dir.IsDir() {
			continue
		}

		localDir := Directory{
			Path:   d.Path + "/" + dir.Name(),
			Fs:     d.Fs,
		}
		localDirs = append(localDirs, adapter.Directory(&localDir))
	}
	return localDirs, nil
}

func (d *Directory) String() string {
	return d.GetPath()
}