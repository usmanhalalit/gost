package adapter

import (
	"io"
	"time"
)


type Filesystem interface {
	GetClient() interface{}
	GetConfig() interface{}
}

type Directory interface {
	File(path string) File
	Files() ([]File, error)
	Directory(path string) Directory
	Filesystem() Filesystem
	GetPath() string
	//Directories() ([]Directory, error)
	//Info()
}

type File interface {
	ReadString() (string, error)
	WriteString(text string) error
	Delete() error
	Exist() bool
	Info() (FileInfo, error)
	Directory() Directory
	GetPath() string
	Filesystem() Filesystem
	io.ReadWriter
}

type FileInfo struct {
	Size int64
	LastModified time.Time
}
