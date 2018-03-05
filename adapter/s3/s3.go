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

type s3adapter struct{
	Service *s3.S3
	Config S3config
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

	return s3adapter{
		Service: s3Service,
		Config: c,
	}
}

func (ad s3adapter) GetString(filename string) (string, error) {
	input := ad.getObjectInput(filename)
	r, err := ad.Service.GetObject(input)
	if err != nil { return "", err}
	text, err := ioutil.ReadAll(r.Body)
	return string(text), err
}

func (ad s3adapter) PutString(filename string, text string) (interface{}, error) {
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader([]byte(text)),
		Bucket: aws.String(ad.Config.Bucket),
		Key:    aws.String(filename),
	}

	r, err := ad.Service.PutObject(input)
	return r, err
}

func (ad s3adapter) Delete(filename string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(ad.Config.Bucket),
		Key:    aws.String(filename),
	}

	_, err := ad.Service.DeleteObject(input)
	return err
}

func (ad s3adapter) GetSignedUrl(filename string, ttl time.Duration) (string, error) {
	req, _ := ad.Service.GetObjectRequest(ad.getObjectInput(filename))
	return req.Presign(ttl)
}

func (ad s3adapter) getObjectInput(filename string) *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(ad.Config.Bucket),
		Key:    aws.String(filename),
	}
}

func (ad s3adapter) Exist(filename string) bool  {
	_, err := ad.Service.GetObject(ad.getObjectInput(filename))
	return err == nil
}

func (ad s3adapter) Info(filename string)  {

}
