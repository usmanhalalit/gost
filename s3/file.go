package s3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

type File struct {
	Path   string
	Fs     *Filesystem
	reader io.ReadCloser
}

func (f *File) Directory() gost.Directory {
	return &Directory{
		Path: filepath.Dir(f.GetPath()),
		Fs: f.Fs,
	}
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) ReadString() (string, error) {
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func (f *File) WriteString(s string) error {
	b := []byte(s)
	n, err := f.Write(b)
	if n != len(b) {
		return errors.New(fmt.Sprintf("Wrote %v bytes from given %v bytes", n, len(b)))
	}
	return err
}

func (f *File) Delete() error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(f.Fs.Config.Bucket),
		Key:    aws.String(f.Path),
	}

	_, err := f.Fs.Service.DeleteObject(input)
	return err
}

func (f *File) Exist() bool  {
	_, err := f.Fs.Service.GetObject(f.getObjectInput())
	return err == nil
}

func (f *File) Stat() (gost.FileInfo, error) {
	info := gost.FileInfo{}

	file, err := f.Fs.Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(f.Fs.Config.Bucket),
		Key:    aws.String(f.Path),
	})
	if err != nil {
		return info, err
	}

	info.Size = *file.ContentLength
	info.LastModified = *file.LastModified

	return info, nil
}

func (f *File) Write(p []byte) (n int, err error) {
	reader := bytes.NewReader(p)
	input := &s3.PutObjectInput{
		Body:   reader,
		Bucket: aws.String(f.Fs.Config.Bucket),
		Key:    aws.String(f.Path),
	}
	_, err = f.Fs.Service.PutObject(input)
	bytesWritten := len(p) - reader.Len()
	// TODO follow rules on io.Writer
	return bytesWritten, err
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		input := f.getObjectInput()
		r, err := f.Fs.Service.GetObject(input)
		if err != nil { return 0, err }
		f.reader = r.Body
	}

	return f.reader.Read(p)
}

func (f *File) Close() error {
	return f.reader.Close()
}

func (f *File) GetSignedUrl(ttl time.Duration) (string, error) {
	req, _ := f.Fs.Service.GetObjectRequest(f.getObjectInput())
	return req.Presign(ttl)
}

func (f *File) getObjectInput() *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(f.Fs.Config.Bucket),
		Key:    aws.String(f.Path),
	}
}

func (f *File) Copy(newName string) error {
	return f.CopyTo(f.Directory(), newName)
}

func (f *File) CopyTo(dir gost.Directory, newName ...string) error {
	var filename string
	if len(newName) > 0 {
		filename = newName[0]
	} else {
		_, filename = filepath.Split(f.GetPath())
	}

	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	file := dir.File(filename)
	n, err := file.Write(content)
	if err != nil {
		return err
	}
	if n != len(content) {
		return errors.New("couldn't copy full file")
	}

	return nil
}

func (f *File) String() string {
	return f.GetPath()
}

