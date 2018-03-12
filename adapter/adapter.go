package adapter

import "time"


type Filesystem interface {
	GetClient() interface{}
	GetConfig() interface{}
}

type Directory interface {
	File(path string) File
	Files() ([]File, error)
	Directory(path string) Directory
	Filesystem() Filesystem
	//Directories() ([]Directory, error)
	//Info()
}

type File interface {
	//Get() ([]byte, error)
	GetString() (string, error)
	//Put(text []byte) (interface{}, error)
	PutString(text string) (interface{}, error)
	Delete() error
	Exist() bool
	Info() (FileInfo, error)
	//Directory() *Directory
	Filesystem() Filesystem
}

type FileInfo struct {
	Size int64
	LastModified time.Time
}
