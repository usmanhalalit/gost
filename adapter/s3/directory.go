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
	path = strings.Trim(path, "/")
	return &S3directory{
		Path: path,
		Fs: ad.Fs,
	}
}

func (ad *S3directory) Files() ([]adapter.File, error) {
	var delimiter *string
	if ad.Path == "" {
		delimiter = aws.String("/")
	} else {
		delimiter = aws.String(ad.Path)
	}

	files, err := ad.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(ad.Fs.Config.Bucket),
		Prefix:    aws.String(ad.Path),
		Delimiter: delimiter,
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

func (ad *S3directory) Directories() ([]adapter.Directory, error) {
	files, err := ad.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(ad.Fs.Config.Bucket),
		Prefix: aws.String(ad.Path),
	})

	if err != nil { return nil, err }
	var s3Directories []adapter.Directory
	addedDirs := make(map[string]bool)

	minNoOfSlash := 2

	if ad.Path != "" {
		slashesInPath := len(strings.Split(ad.Path, "/"))
		if slashesInPath > 0 {
			minNoOfSlash += slashesInPath
		}
	}

	for i := range files.Contents {
		filename := *files.Contents[i].Key
		parts := strings.Split(filename, "/")
		if len(parts) < minNoOfSlash {
			continue
		}

		dir := parts[0]
		if _, ok := addedDirs[dir]; ok {
			continue
		}

		parts = parts[:len(parts)-1]
		s3directory := S3directory {
			Path: strings.Join(parts, "/"),
			Fs:   ad.Fs,
		}
		// TODO may need a fix
		addedDirs[dir] = true
		s3Directories = append(s3Directories, adapter.Directory(&s3directory))
	}
	return s3Directories, nil
}

func (ad *S3directory) String() string {
	return ad.GetPath()
}

func (ad *S3directory) getObjectInput() *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(ad.Fs.Config.Bucket),
		Key:    aws.String(ad.Path),
	}
}

func (ad *S3directory) Exist() bool  {
	list, err := ad.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(ad.Fs.Config.Bucket),
		Prefix:    aws.String(ad.Path + "/"),
		MaxKeys: aws.Int64(1),
	})

	return err == nil && len(list.Contents) > 0
}

func (ad *S3directory) Create() error {
	reader := strings.NewReader("")
	input := &s3.PutObjectInput{
		Body:   reader,
		Bucket: aws.String(ad.Fs.Config.Bucket),
		Key:    aws.String(ad.Path + "/"),
	}
	_, err := ad.Fs.Service.PutObject(input)
	return err
}

func (ad *S3directory) Delete() error {
	files, err := ad.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(ad.Fs.Config.Bucket),
		Prefix:    aws.String(ad.Path),
	})

	if err != nil {
		return err
	}

	for i := range files.Contents {
		doi := &s3.DeleteObjectInput{
			Bucket:    aws.String(ad.Fs.Config.Bucket),
			Key: files.Contents[i].Key,
		}
		_, err = ad.Fs.Service.DeleteObject(doi)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ad *S3directory) Stat() (adapter.FileInfo, error) {
	panic("Stat is not available on S3 directory")
}