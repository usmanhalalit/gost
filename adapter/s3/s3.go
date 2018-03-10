package s3

import (
	"github.com/usmanhalalit/gost/adapter"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"bytes"
	"io/ioutil"
	"strings"
	"time"
)

type S3filesystem struct{
	S3directory
	Service *s3.S3
	Config S3config
}

type S3file struct {
	Path string
	Filesystem *S3filesystem
}

type S3directory struct {
	Path string
	Filesystem *S3filesystem
}

type S3config struct {
	Id string
	Secret string
	Token string
	Region string
	Bucket string
}
var s3Service *s3.S3

func NewS3Adapter(c S3config) adapter.Directory {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.Id, c.Secret, c.Token),
	})
	

	// Create S3 service client with a specific Region.
	s3Service = s3.New(sess)

	fs := S3filesystem{
		Service: s3Service,
		Config: c,
	}
	return S3directory{
		Filesystem: &fs,
		Path: "",
	}
}

func (ad S3directory) GetClient() interface{} {
	return ad.Filesystem.Service
}

func (ad S3directory) GetConfig() interface{} {
	return ad.Filesystem.Config
}

func (ad S3directory) File(path string) adapter.File {
	return &S3file{
		Path: path,
		Filesystem: ad.Filesystem,
	}
}

func (ad S3directory) Directory(path string) adapter.Directory {
	path = ad.Path + "/" + path
	path = strings.Trim(path, "/")
	return S3directory{
		Path: path,
		Filesystem: ad.Filesystem,
	}
}

func (ad S3directory) Files() ([]adapter.File, error) {
	files, err := ad.Filesystem.Service.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(ad.Filesystem.Config.Bucket),
		Prefix: aws.String(ad.Path),
	})
	if err != nil { return nil, err }
	var s3files []adapter.File
	for i := range files.Contents {
		s3file := S3file{
			Path: *files.Contents[i].Key,
			Filesystem: ad.Filesystem,
		}
		s3files = append(s3files, adapter.File(s3file))
	}
	return s3files, nil
}

//func (ad S3file) Directory() *adapter.Directory {
//	return S3filesystem{}
//}

func (f S3file) GetString() (string, error) {
	input := f.getObjectInput()
	r, err := f.Filesystem.Service.GetObject(input)
	if err != nil { return "", err}
	text, err := ioutil.ReadAll(r.Body)
	return string(text), err
}

func (f S3file) PutString(text string) (interface{}, error) {
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader([]byte(text)),
		Bucket: aws.String(f.Filesystem.Config.Bucket),
		Key:    aws.String(f.Path),
	}

	r, err := f.Filesystem.Service.PutObject(input)
	return r, err
}

func (f S3file) Delete() error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(f.Filesystem.Config.Bucket),
		Key:    aws.String(f.Path),
	}

	_, err := f.Filesystem.Service.DeleteObject(input)
	return err
}

func (f S3file) Exist() bool  {
	_, err := f.Filesystem.Service.GetObject(f.getObjectInput())
	return err == nil
}

func (f S3file) Info() (adapter.FileInfo, error) {
	info := adapter.FileInfo{}

	file, err := f.Filesystem.Service.GetObject(f.getObjectInput())
	if err != nil {
		return info, err
	}

	info.Size = *file.ContentLength
	info.LastModified = *file.LastModified

	return info, nil
}

func (f S3file) GetSignedUrl(ttl time.Duration) (string, error) {
	req, _ := f.Filesystem.Service.GetObjectRequest(f.getObjectInput())
	return req.Presign(ttl)
}

func (f S3file) getObjectInput() *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(f.Filesystem.Config.Bucket),
		Key:    aws.String(f.Path),
	}
}