package adapter

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
	fmt.Stringer
	// TODO Delete() error
	// TODO Create() error
	 Directories() ([]Directory, error)
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
	fmt.Stringer
	Filesystem() Filesystem
	// TODO io.Closer
	io.ReadWriter
}

// TODO make it Go os.fileinfo compatible
type FileInfo struct {
	Size int64
	LastModified time.Time
}
