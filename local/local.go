package local

import (
	"errors"
	"github.com/usmanhalalit/gost"
)

type LocalFilesystem struct {
	Service interface{}
	Config  Config
}

type Config struct {
	BasePath string
}

func New(c Config) (gost.Directory, error) {
	fs := LocalFilesystem{
		Service: nil,
		Config:  c,
	}

	rootDir := Directory{
		Object{
			Fs:   &fs,
			Path: c.BasePath,
		},
	}

	// Checking if we can read from the directory
	if _, err := rootDir.Stat(); err != nil {
		return nil, errors.New("couldn't read, either directory or it's permission is invalid")
	}

	return &rootDir, nil
}
