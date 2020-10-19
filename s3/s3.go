package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/usmanhalalit/gost"
)

type Filesystem struct {
	Service s3iface.S3API
	Config  Config
}

type Config struct {
	Id     string
	Secret string
	Token  string
	Region string
	Bucket string
}

var globalService s3iface.S3API

func New(c Config, service s3iface.S3API) (gost.Directory, error) {
	// service always takes priority over globalService
	if globalService != nil && service == nil {
		service = globalService
	}

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
	rootDir := &Directory{
		Fs:   &fs,
		Path: "",
	}

	// Checking if we can read from the directory
	if _, err := rootDir.Files(); err != nil {
		return nil, errors.New("couldn't read from S3, credentials could be invalid")
	}

	return rootDir, nil
}

// Deprecated: SetService will be removed in the next major release
func SetService(s s3iface.S3API) {
	globalService = s
}
