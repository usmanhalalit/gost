package adapter

import "time"

type Filesystem interface {
	GetString(filename string) (string, error)
	Delete(filename string) error
	PutString(filename string, text string) (interface{}, error)
	GetSignedUrl(filename string, ttl time.Duration) (string, error)
	Exist(filename string) bool
	Info(filename string)
}
