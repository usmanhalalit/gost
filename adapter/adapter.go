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
	// TODO files() returns recursive files
	Files() ([]File, error)
	Directory(path string) Directory
	Filesystem() Filesystem
	GetPath() string
	// TODO Directories() ([]Directory, error)
	// TODO Stat()
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
	io.ReadWriter
}

// TODO make it Go os.fileinfo compatible
type FileInfo struct {
	Size int64
	LastModified time.Time
}
