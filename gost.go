package gost

import (
	"fmt"
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
	Delete() error
	Exist() bool
	Create() error
	Directories() ([]Directory, error)
	Stat() (FileInfo, error)
	fmt.Stringer
}

type File interface {
	ReadString() (string, error)
	WriteString(text string) error
	Delete() error
	Exist() bool
	Stat() (FileInfo, error)
	Directory() Directory
	GetPath() string
	Filesystem() Filesystem
	io.ReadWriteCloser
	fmt.Stringer
}

// TODO make it Go os.fileinfo compatible
type FileInfo struct {
	Size         int64
	LastModified time.Time
}
