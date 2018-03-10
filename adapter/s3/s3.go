package s3

import (
	"github.com/usmanhalalit/gost/adapter"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"bytes"
	"io/ioutil"
	"time"
)

type S3adapter struct{
	Service *s3.S3
	Config S3config
}

type S3file struct {
	Path string
	Filesystem *S3adapter
}

type S3directory struct {
	Path string
	Filesystem *S3adapter
}

type S3config struct {
	Id string
	Secret string
	Token string
	Region string
	Bucket string
}
var s3Service *s3.S3

func NewS3Adapter(c S3config) adapter.Filesystem {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.Id, c.Secret, c.Token),
	})
	

	// Create S3 service client with a specific Region.
	s3Service = s3.New(sess)

	return S3adapter{
		Service: s3Service,
		Config: c,
	}
}

func (ad S3adapter) GetClient() interface{} {
	return ad.Service
}

func (ad S3adapter) GetConfig() interface{} {
	return ad.Config
}

func (ad S3adapter) File(path string) adapter.File {
	return &S3file{
		Path: path,
		Filesystem: &ad,
	}
}

func (ad S3adapter) Files() ([]adapter.File, error) {
	var s3files []S3file
	for i := 0; i < 10; i++ {
		s3files = append(s3files, S3file{})
	}
	return s3files, nil
}

//func (ad S3file) Directory() *adapter.Directory {
//	return S3adapter{}
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