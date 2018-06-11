// Filesystem abstraction layer for Golang, that works with Local file system
// and Amazon S3 with a unified API. You can even copy-paste files from different sources.
// FTP, Dropbox etc. will follow soon.
//
// Full documentation: https://github.com/usmanhalalit/gost/blob/master/Readme.md
package gost

import (
	"fmt"
	"io"
	"time"
)

type Directory interface {
	// Point to a file in directory
	File(path string) File
	// Get all files from the directory
	Files() ([]File, error)
	// Point to a directory in this directory
	Directory(path string) Directory
	// Get the the directory path
	GetPath() string
	// Delete the entire directory, no undo
	Delete() error
	// Check if directory exists
	Exists() bool
	// Create a directory
	Create() error
	// Get all directories in the directory
	Directories() ([]Directory, error)
	// Get directory metadata
	Stat() (FileInfo, error)
	// Helper to convert directory to a string, provide directory name
	fmt.Stringer
}

type File interface {
	// Get all file content as string
	ReadString() (string, error)
	// Write string to the file
	WriteString(text string) error
	// Delete the file, no undo
	Delete() error
	// Check of file exists
	Exists() bool
	// Get file metadata
	Stat() (FileInfo, error)
	// Get the directory which contains the file
	Directory() Directory
	// Get the file path
	GetPath() string
	// Copy file to another directory
	CopyTo(dir Directory, newName ...string) error
	// Copy file to the current directory
	Copy(newName string) error
	// Read to file, write to file and close the file
	io.ReadWriteCloser
	// Helper to convert file to a string, provide file name
	fmt.Stringer
}

type FileInfo struct {
	Size         int64
	LastModified time.Time
}
