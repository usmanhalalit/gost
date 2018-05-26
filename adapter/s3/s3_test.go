package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/usmanhalalit/gost/adapter"
	"github.com/usmanhalalit/gost/mocks"
	"io"
	"io/ioutil"
	"log"
	"testing"
	"time"
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
		Key:    aws.String("/test.txt"),
		Body: 	bytes.NewReader([]byte("test")),
	}).Return(&s3.PutObjectOutput{}, nil)
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

	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 512 {
		t.Errorf("Wrong byte size on read")
	}
}

func Test_GetString(t *testing.T) {
	_, err := s3fs.File("test.txt").ReadString()
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
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
	if err != nil {
		t.Errorf("Couldn't get info: %v", err)
	}

	if info.Size != 3 {
		t.Errorf("Invalid file size expected %v got %v", 3, info.Size)
	}

	if info.LastModified.Day() != time.Now().Day() {
		t.Errorf("Invalid file size expected %v got %v", time.Now().Day(), info.LastModified.Day())
	}
}

func Test_Delete(t *testing.T) {
	//readErr = errors.New("file does not exist")

	err := s3fs.File("test.txt").Delete()
	if err != nil {
		t.Errorf("Failed deleting: %v", err)
	}
	_, err = s3fs.File("test.txt").ReadString()
}

//func Test_NotExist(t *testing.T)  {
//	if s3fs.File("test.txt").Exist() {
//		t.Errorf("File does exist")
//	}
//}

func Test_Write_In_Sub_Dir(t *testing.T) {
	err := s3fs.File("aDir/aDirSub/subsub.txt").WriteString("abc")
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_Directories(t *testing.T) {
	dirs, _ := s3fs.Directory("aDir").Directories()
	log.Println(dirs)
	//dirs[0].File("subsub.txt").ReadString()
}

func Test_Files_In_Dir(t *testing.T) {
	files, err := s3fs.Directory("/").Files()

	if len(files) < 1 {
		t.Fatalf("Failed listing found %v files", len(files))
	}

	log.Println(files[0].GetPath())

	if err != nil {
		t.Errorf("Failed listing: %v", err)
	}
}

func Test_Create_Dir(t *testing.T) {
	if err := s3fs.Directory("dDir").Create(); err != nil {
		t.Errorf("Couldn't create dir")
	}
}

func Test_Exist_Dir(t *testing.T) {
	if ! s3fs.Directory("dDir").Exist() {
		t.Errorf("Dir doesn't exist")
	}
}

func Test_Delete_Dir(t *testing.T) {
	if err := s3fs.Directory("dDir").Delete(); err != nil {
		t.Errorf("Couldn't delete dir: %v", err)
	}
}

func Test_Dir_Not_Exist(t *testing.T) {
	if s3fs.Directory("dDir").Exist() {
		t.Errorf("Dir exists")
	}
}

