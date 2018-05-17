package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
)

type S3Service interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObjectRequest(input *s3.GetObjectInput) (req *request.Request, output *s3.GetObjectOutput)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
}

type Filesystem struct {
	Service S3Service
	Config  Config
}


type Config struct {
	Id string
	Secret string
	Token string
	Region string
	Bucket string
}

var service S3Service

func New(c Config) adapter.Directory {
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

func (ad *Filesystem) GetClient() interface{} {
	return ad.Service
}

func (ad *Filesystem) GetConfig() interface{} {
	return ad.Config
}

func SetService(s S3Service) {
	service = s
}