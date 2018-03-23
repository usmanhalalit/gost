package local

import (
	"github.com/usmanhalalit/gost/adapter"
)

type LocalFilesystem struct {
	Service interface{}
	Config LocalConfig
}

type LocalConfig struct {
	BasePath string
}

func NewLocalAdapter(c LocalConfig) adapter.Directory {
	fs := LocalFilesystem{
		Service: nil,
		Config: c,
	}
	return &Directory{
		Fs: &fs,
		Path: c.BasePath,
	}
}

func (fs *LocalFilesystem) GetClient() interface{} {
	return fs.Service
}

func (fs *LocalFilesystem) GetConfig() interface{} {
	return fs.Config
}