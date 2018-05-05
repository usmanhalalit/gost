package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
)

type Filesystem struct {
	Service *s3.S3
	Config  Config
}


type Config struct {
	Id string
	Secret string
	Token string
	Region string
	Bucket string
}

var service *s3.S3

func New(c Config) adapter.Directory {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.Id, c.Secret, c.Token),
	})
	

	// Create S3 service client with a specific Region.
	service = s3.New(sess)

	fs := Filesystem{
		Service: service,
		Config:  c,
	}
	return &Directory{
		Fs: &fs,
		Path: "",
	}
}

func (ad *Filesystem) GetClient() interface{} {
	return ad.Service
}

func (ad *Filesystem) GetConfig() interface{} {
	return ad.Config
}