package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/usmanhalalit/gost"
	"github.com/usmanhalalit/gost/mocks"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

var s3fs gost.Directory
var s3mock mocks.S3API

func init() {
	s3mock = mocks.S3API{}
	SetService(&s3mock)

	setMockExpectations()

	s3fs, _ = New(Config{
		Bucket: "fake",
		Region: "eu-west-1",
		Id:     "aws_access_id",
		Secret: "aws_secret_id",
	})

}

func setMockExpectations() {
	rwc := mocks.ReadWriteCloseSeeker{}
	rwc.On("Read", make([]byte, 512)).Return(512, io.EOF)
	s3mock.On("GetObject", &s3.GetObjectInput{
		Bucket: aws.String("fake"),
		Key:    aws.String("/test.txt"),
	}).Return(&s3.GetObjectOutput{
		Body: rwc,
	}, nil)
	s3mock.On("HeadObject", &s3.HeadObjectInput{
		Bucket: aws.String("fake"),
		Key:    aws.String("/test.txt"),
	}).Return(&s3.HeadObjectOutput{
		ContentLength: aws.Int64(3),
		LastModified:  aws.Time(time.Now()),
	}, nil)
	s3mock.On("DeleteObject", &s3.DeleteObjectInput{
		Bucket: aws.String("fake"),
		Key:    aws.String("/test.txt"),
	}).Return(&s3.DeleteObjectOutput{}, nil)
	s3mock.On("PutObject", &s3.PutObjectInput{
		Bucket: aws.String("fake"),
		Key:    aws.String("fake_new_dir/"),
		Body:   strings.NewReader(""),
	}).Return(&s3.PutObjectOutput{}, nil)
	s3mock.On("PutObject", &s3.PutObjectInput{
		Bucket: aws.String("fake"),
		Key:    aws.String("/aDir/aDirSub/subsub.txt"),
		Body:   bytes.NewReader([]byte("test")),
	}).Return(&s3.PutObjectOutput{}, nil)
	s3mock.On("PutObject", &s3.PutObjectInput{
		Bucket: aws.String("fake"),
		Key:    aws.String("/test.txt"),
		Body:   bytes.NewReader([]byte("test")),
	}).Return(&s3.PutObjectOutput{}, nil)
	keys := []*s3.Object{
		{Key: aws.String("test.txt")},
		{Key: aws.String("aDir/test.txt")},
		{Key: aws.String("aDir/subdir/test.txt")},
		{Key: aws.String("bDir/subdir/test.txt")},
	}
	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix: aws.String("fake_new_dir"),
	}).Return(&s3.ListObjectsOutput{
		Contents: []*s3.Object{
			{Key: aws.String("/test.txt")},
		},
	}, nil)
	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket:  aws.String("fake"),
		MaxKeys: aws.Int64(1),
		Prefix:  aws.String("fake_new_dir/"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)
	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket:    aws.String("fake"),
		Prefix:    aws.String("aDir"),
		Delimiter: aws.String("aDir"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)
	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix: aws.String(""),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)
	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix: aws.String("aDir"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)
	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix: aws.String(""),
		Delimiter: aws.String("/"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)
}

func TestWrite(t *testing.T) {
	f := s3fs.File("test.txt")
	n, err := f.Write([]byte("test"))

	if n != 0 {
		t.Errorf("Failed writing as io.Writer wrote %v bytes should be %v bytes", n, 0)
	}

	if err != nil {
		t.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	f := s3fs.File("test.txt")
	b, err := ioutil.ReadAll(f)

	assert.NoError(t, err)
	if len(b) != 512 {
		t.Errorf("Wrong byte size on read")
	}
}

func TestGetString(t *testing.T) {
	_, err := s3fs.File("test.txt").ReadString()
	assert.NoError(t, err)
}

//func Test_GetSignedUrl(t *testing.T) {
//	f := s3fs.File("test.txt")
//	_, err := f.(*File).GetSignedUrl(time.Minute * 1)
//	if err != nil {
//		t.Errorf("Failed write: %v", err)
//	}
//}

func TestExist(t *testing.T) {
	if !s3fs.File("test.txt").Exists() {
		t.Errorf("File doesn't exist")
	}
}

func TestStat(t *testing.T) {
	info, err := s3fs.File("test.txt").Stat()
	assert.NoError(t, err)

	if info.Size != 3 {
		t.Errorf("Invalid file size expected %v got %v", 3, info.Size)
	}

	if info.LastModified.Day() != time.Now().Day() {
		t.Errorf("Invalid file size expected %v got %v", time.Now().Day(), info.LastModified.Day())
	}
}

func TestDelete(t *testing.T) {
	err := s3fs.File("test.txt").Delete()
	assert.NoError(t, err)
}

//func Test_NotExist(t *testing.T)  {
//	if s3fs.File("test.txt").Exists() {
//		t.Errorf("File does exist")
//	}
//}

func TestWriteInSubDir(t *testing.T) {
	_, err := s3fs.File("aDir/aDirSub/subsub.txt").Write([]byte("test"))
	assert.NoError(t, err)
}

func TestDirectories(t *testing.T) {
	dirs, err := s3fs.Directory("aDir").Directories()
	dirsStr := fmt.Sprintf("%v", dirs)
	expectedDirStr := "[aDir/subdir bDir/subdir]"
	assert.Equal(t, expectedDirStr, dirsStr, "did not get expected directories, expected: %v got %v", expectedDirStr, dirsStr)
	assert.NoError(t, err)

	dirs, err = s3fs.Directories()
	dirsStr = fmt.Sprintf("%v", dirs)
	expectedDirStr = "[aDir bDir/subdir]"
	assert.Equal(t, expectedDirStr, dirsStr, "did not get expected directories, expected: %v got %v", expectedDirStr, dirsStr)
	assert.NoError(t, err)
}

func TestFilesInDir(t *testing.T) {
	files, err := s3fs.Directory("aDir").Files()

	assert.Condition(t, func() bool {
		return len(files) >= 1
	}, "Failed listing")
	assert.Equal(t, files[0].GetPath(), "test.txt")
	assert.NoError(t, err)
}

func TestCreateDir(t *testing.T) {
	assert.NoError(t, s3fs.Directory("fake_new_dir").Create())
}

func TestExistDir(t *testing.T) {
	if !s3fs.Directory("fake_new_dir").Exists() {
		t.Errorf("Dir doesn't exist")
	}
}

func TestDeleteDir(t *testing.T) {
	assert.NoError(t, s3fs.Directory("fake_new_dir").Delete())
}
