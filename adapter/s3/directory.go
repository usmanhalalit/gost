package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
	"strings"
)

type S3directory struct {
	Path string
	Fs *S3filesystem
}

func (ad *S3directory) Filesystem() adapter.Filesystem {
	return ad.Fs
}

func (ad *S3directory) File(path string) adapter.File {
	return &S3file{
		Path:   ad.Path + "/" + path,
		Fs:     ad.Fs,
		reader: nil,
	}
}

func (ad *S3directory) GetPath() string {
	return ad.Path
}

func (ad *S3directory) Directory(path string) adapter.Directory {
	path = ad.Path + "/" + path
	path = strings.TrimRight(path, "/")
	return &S3directory{
		Path: path,
		Fs: ad.Fs,
	}
}

func (ad *S3directory) Files() ([]adapter.File, error) {
	files, err := ad.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(ad.Fs.Config.Bucket),
		Prefix: aws.String(ad.Path),
	})
	if err != nil { return nil, err }
	var s3files []adapter.File
	for i := range files.Contents {
		s3file := S3file{
			Path: *files.Contents[i].Key,
			Fs: ad.Fs,
			reader: nil,
		}
		s3files = append(s3files, adapter.File(&s3file))
	}
	return s3files, nil
}
