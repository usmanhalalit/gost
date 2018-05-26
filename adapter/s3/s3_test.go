package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
	"github.com/usmanhalalit/gost/mocks"
	"io"
	"io/ioutil"
	"testing"
	"time"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
)

var s3fs adapter.Directory
var s3mock mocks.S3API

func init() {
	s3mock = mocks.S3API{}
	SetService(&s3mock)
	s3fs = New(Config{
		Bucket: "fake",
	})

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
		Prefix:aws.String("fake_new_dir"),
	}).Return(&s3.ListObjectsOutput{
		Contents: []*s3.Object{
			{Key: aws.String("/test.txt")},
		},
	}, nil)

	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		MaxKeys: aws.Int64(1),
		Prefix:aws.String("fake_new_dir/"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)

	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix:    aws.String("aDir"),
		Delimiter:    aws.String("aDir"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)

	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix:    aws.String(""),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)

	s3mock.On("ListObjects", &s3.ListObjectsInput{
		Bucket: aws.String("fake"),
		Prefix:    aws.String("aDir"),
	}).Return(&s3.ListObjectsOutput{
		Contents: keys,
	}, nil)

}

func Test_Write(t *testing.T) {
	f := s3fs.File("test.txt")
	n, err := f.Write([]byte("test"))

	if n != 0 {
		t.Errorf("Failed writing as io.Writer wrote %v bytes should be %v bytes", n, 0)
	}

	if err != nil {
		t.Fatal(err)
	}
}

func Test_Read(t *testing.T) {
	f := s3fs.File("test.txt")
	b, err := ioutil.ReadAll(f)

	assert.NoError(t, err)
	if len(b) != 512 {
		t.Errorf("Wrong byte size on read")
	}
}

func Test_GetString(t *testing.T) {
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

func Test_Exist(t *testing.T) {
	if ! s3fs.File("test.txt").Exist() {
		t.Errorf("File doesn't exist")
	}
}

func Test_Stat(t *testing.T) {
	info, err := s3fs.File("test.txt").Stat()
	assert.NoError(t, err)

	if info.Size != 3 {
		t.Errorf("Invalid file size expected %v got %v", 3, info.Size)
	}

	if info.LastModified.Day() != time.Now().Day() {
		t.Errorf("Invalid file size expected %v got %v", time.Now().Day(), info.LastModified.Day())
	}
}

func Test_Delete(t *testing.T) {
	err := s3fs.File("test.txt").Delete()
	assert.NoError(t, err)
}

//func Test_NotExist(t *testing.T)  {
//	if s3fs.File("test.txt").Exist() {
//		t.Errorf("File does exist")
//	}
//}

func Test_Write_In_Sub_Dir(t *testing.T) {
	_, err := s3fs.File("aDir/aDirSub/subsub.txt").Write([]byte("test"))
	assert.NoError(t, err)
}

func Test_Directories(t *testing.T) {
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


func Test_Files_In_Dir(t *testing.T) {
	files, err := s3fs.Directory("aDir").Files()

	assert.Condition(t, func() bool {
		return len(files) >= 1
	}, "Failed listing")
	assert.Equal(t, files[0].GetPath(), "test.txt")
	assert.NoError(t, err)
}

func Test_Create_Dir(t *testing.T) {
	assert.NoError(t, s3fs.Directory("fake_new_dir").Create())
}

func Test_Exist_Dir(t *testing.T) {
	if ! s3fs.Directory("fake_new_dir").Exist() {
		t.Errorf("Dir doesn't exist")
	}
}

func Test_Delete_Dir(t *testing.T) {
	assert.NoError(t, s3fs.Directory("fake_new_dir").Delete())
}
