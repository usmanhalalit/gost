package adapter

import "time"

type Filesystem interface {
	GetString(filename string) (string, error)
	Delete(filename string) error
	PutString(filename string, text string) (interface{}, error)
	Exist(filename string) bool
	Info(filename string) (FileInfo, error)
}

type FileInfo struct {
	Size int64
	LastModified time.Time
}
