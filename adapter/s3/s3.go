package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
)

type S3filesystem struct {
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
	return &S3directory{
		Fs: &fs,
		Path: "",
	}
}

func (ad *S3filesystem) GetClient() interface{} {
	return ad.Service
}

func (ad *S3filesystem) GetConfig() interface{} {
	return ad.Config
}