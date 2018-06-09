package local

import (
	"github.com/usmanhalalit/gost"
	"io/ioutil"
	"os"
	"strings"
)

type Directory struct {
	Object
}

func (d *Directory) File(path string) gost.File {
	return &File{
		Object: Object{
			Path: d.Path + "/" + path,
			Fs:   d.Fs,
		},
		reader: nil,
	}
}

func (d *Directory) Directory(path string) gost.Directory {
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

func (d *Directory) Files() ([]gost.File, error) {
	files, err := ioutil.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}
	var localFiles []gost.File
	for _, file := range files {
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
		localFiles = append(localFiles, gost.File(&localFile))
	}
	return localFiles, nil
}

func (d *Directory) Directories() ([]gost.Directory, error) {
	files, err := ioutil.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}
	var localDirs []gost.Directory
	for _, dir := range files {
		if !dir.IsDir() {
			continue
		}

		localDir := Directory{
			Object: Object{
				Path: d.Path + "/" + dir.Name(),
				Fs:   d.Fs,
			},
		}
		localDirs = append(localDirs, gost.Directory(&localDir))
	}
	return localDirs, nil
}
