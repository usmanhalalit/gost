package local

import (
	"github.com/usmanhalalit/gost/adapter"
	"io/ioutil"
	"os"
	"strings"
)

type Directory struct {
	Object
}

func (d *Directory) File(path string) adapter.File {
	return &File{
		Object: Object{
			Path: d.Path + "/" + path,
			Fs:   d.Fs,
		},
		reader: nil,
	}
}

func (d *Directory) Directory(path string) adapter.Directory {
	path = d.Path + "/" + path
	path = strings.TrimRight(path, "/")
	return &Directory{
		Object{
			Path: path,
			Fs:   d.Fs,
		},
	}
}

func (d *Directory) Create() error {
	return os.Mkdir(d.Path, 644)
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
			Object: Object{
				Path: d.Path + "/" + file.Name(),
				Fs:   d.Fs,
			},
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
			Object: Object{
				Path: d.Path + "/" + dir.Name(),
				Fs:   d.Fs,
			},
		}
		localDirs = append(localDirs, adapter.Directory(&localDir))
	}
	return localDirs, nil
}