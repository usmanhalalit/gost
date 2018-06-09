package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Filesystem struct {
	Service s3iface.S3API
	Config  Config
}


type Config struct {
	Id string
	Secret string
	Token string
	Region string
	Bucket string
}

var service s3iface.S3API

func New(c Config) gost.Directory {
	if service == nil {
		sess, _ := session.NewSession(&aws.Config{
			Region:      aws.String(c.Region),
			Credentials: credentials.NewStaticCredentials(c.Id, c.Secret, c.Token),
		})
		service = s3.New(sess)
	}

	fs := Filesystem{
		Service: service,
		Config:  c,
	}
	return &Directory{
		Fs: &fs,
		Path: "",
	}
}

func SetService(s s3iface.S3API) {
	service = s
}