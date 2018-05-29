package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost"
	"strings"
)

type Directory struct {
	Path string
	Fs *Filesystem
}

func (d *Directory) Filesystem() gost.Filesystem {
	return d.Fs
}

func (d *Directory) File(path string) gost.File {
	return &File{
		Path:   d.Path + "/" + path,
		Fs:     d.Fs,
		reader: nil,
	}
}

func (d *Directory) GetPath() string {
	return d.Path
}

func (d *Directory) Directory(path string) gost.Directory {
	path = d.Path + "/" + path
	path = strings.Trim(path, "/")
	return &Directory{
		Path: path,
		Fs:   d.Fs,
	}
}

func (d *Directory) Files() ([]gost.File, error) {
	var delimiter *string
	if d.Path == "" {
		delimiter = aws.String("/")
	} else {
		delimiter = aws.String(d.Path)
	}

	files, err := d.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(d.Fs.Config.Bucket),
		Prefix:    aws.String(d.Path),
		Delimiter: delimiter,
	})

	if err != nil { return nil, err }
	var s3files []gost.File
	for i := range files.Contents {
		s3file := File{
			Path:   *files.Contents[i].Key,
			Fs:     d.Fs,
			reader: nil,
		}
		s3files = append(s3files, gost.File(&s3file))
	}
	return s3files, nil
}

func (d *Directory) Directories() ([]gost.Directory, error) {
	files, err := d.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(d.Fs.Config.Bucket),
		Prefix: aws.String(d.Path),
	})

	if err != nil { return nil, err }
	var s3Directories []gost.Directory
	addedDirs := make(map[string]bool)

	minNoOfSlash := 2

	if d.Path != "" {
		slashesInPath := len(strings.Split(d.Path, "/"))
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
		s3directory := Directory{
			Path: strings.Join(parts, "/"),
			Fs:   d.Fs,
		}
		// TODO may need a fix
		addedDirs[dir] = true
		s3Directories = append(s3Directories, gost.Directory(&s3directory))
	}
	return s3Directories, nil
}

func (d *Directory) String() string {
	return d.GetPath()
}

func (d *Directory) getObjectInput() *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(d.Fs.Config.Bucket),
		Key:    aws.String(d.Path),
	}
}

func (d *Directory) Exist() bool  {
	list, err := d.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(d.Fs.Config.Bucket),
		Prefix:    aws.String(d.Path + "/"),
		MaxKeys: aws.Int64(1),
	})

	return err == nil && len(list.Contents) > 0
}

func (d *Directory) Create() error {
	reader := strings.NewReader("")
	input := &s3.PutObjectInput{
		Body:   reader,
		Bucket: aws.String(d.Fs.Config.Bucket),
		Key:    aws.String(d.Path + "/"),
	}
	_, err := d.Fs.Service.PutObject(input)
	return err
}

func (d *Directory) Delete() error {
	files, err := d.Fs.Service.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(d.Fs.Config.Bucket),
		Prefix:    aws.String(d.Path),
	})

	if err != nil {
		return err
	}

	for i := range files.Contents {
		doi := &s3.DeleteObjectInput{
			Bucket:    aws.String(d.Fs.Config.Bucket),
			Key: files.Contents[i].Key,
		}
		_, err = d.Fs.Service.DeleteObject(doi)

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Directory) Stat() (gost.FileInfo, error) {
	panic("Stat is not available on S3 directory")
}