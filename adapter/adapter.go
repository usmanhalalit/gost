package adapter

import "time"


type Filesystem interface {
	GetClient() interface{}
	GetConfig() interface{}
}

type Directory interface {
	Filesystem
	File(path string) File
	Files() ([]File, error)
	Directory(path string) Directory
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
}

type FileInfo struct {
	Size int64
	LastModified time.Time
}
