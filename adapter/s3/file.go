package s3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

type S3file struct {
	Path   string
	Fs     *S3filesystem
	reader io.Reader
	writer io.WriteCloser
}

func (f *S3file) Directory() adapter.Directory {
	return &S3directory{
		Path: filepath.Dir(f.GetPath()),
		Fs: f.Fs,
	}
}

func (f *S3file) Filesystem() adapter.Filesystem {
	return f.Fs
}

func (f *S3file) GetPath() string {
	return f.Path
}

func (f *S3file) ReadString() (string, error) {
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func (f *S3file) WriteString(s string) error {
	b := []byte(s)
	n, err := f.Write(b)
	if n != len(b) {
		return errors.New(fmt.Sprintf("Wrote %v bytes from given %v bytes", n, len(b)))
	}
	return err
}

func (f *S3file) Delete() error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(f.Fs.Config.Bucket),
		Key:    aws.String(f.Path),
	}

	_, err := f.Fs.Service.DeleteObject(input)
	return err
}

func (f *S3file) Exist() bool  {
	_, err := f.Fs.Service.GetObject(f.getObjectInput())
	return err == nil
}

func (f *S3file) Info() (adapter.FileInfo, error) {
	info := adapter.FileInfo{}

	file, err := f.Fs.Service.GetObject(f.getObjectInput())
	if err != nil {
		return info, err
	}

	info.Size = *file.ContentLength
	info.LastModified = *file.LastModified

	return info, nil
}

func (f *S3file) Write(p []byte) (n int, err error) {
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

func (f *S3file) Read(p []byte) (n int, err error) {
	if f.reader == nil {
		input := f.getObjectInput()
		r, err := f.Fs.Service.GetObject(input)
		if err != nil { return 0, err }
		f.reader = r.Body
	}

	return f.reader.Read(p)
}

func (f *S3file) GetSignedUrl(ttl time.Duration) (string, error) {
	req, _ := f.Fs.Service.GetObjectRequest(f.getObjectInput())
	return req.Presign(ttl)
}

func (f *S3file) getObjectInput() *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(f.Fs.Config.Bucket),
		Key:    aws.String(f.Path),
	}
}

