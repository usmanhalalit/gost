package s3

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
	"github.com/usmanhalalit/gost/adapter"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var s3fs adapter.Directory

type mockS3 struct {
	mock.Mock
}

type ReadWriteCloser struct {
}

func (ReadWriteCloser) Read(p []byte) (n int, err error) {
	return 1, nil
}

func (ReadWriteCloser) Write(p []byte) (n int, err error) {
	return 1, nil
}

func (ReadWriteCloser) Close() error {
	return nil
}

func (m mockS3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}
func (mockS3) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) { return nil, nil}
func (mockS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) { return nil, nil}
func (mockS3) GetObjectRequest(input *s3.GetObjectInput) (req *request.Request, output *s3.GetObjectOutput) { return nil, nil}
func (mockS3) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) { return nil, nil}
func (mockS3) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) { return nil, nil}

func init() {
	s3fs = New(Config{
		Id: "AKIAJBRFB4PEZIKTETJQ",
		Secret: "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",
		Region: "us-east-1",
		Bucket: "usman-gost",
	})
}

func Test_New(t *testing.T) {
	err := s3fs.File("test.txt").WriteString("abc")
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_Write_In_Sub_Dir(t *testing.T) {
	err := s3fs.File("aDir/aDirSub/subsub.txt").WriteString("abc")
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_Files(t *testing.T) {
	files, _ := s3fs.Files()
	log.Println(files)
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

func Test_Write(t *testing.T) {
	f := s3fs.File("firas.jpg")
	firas, err := os.Open("../../storage/firas.jpg")
	if err != nil {
		t.Fatal(err)
	}
	fi, _ := firas.Stat()
	size := int(fi.Size())
	firasB := make([]byte, size)
	_, err = firas.Read(firasB)
	if err != nil {
		t.Fatal(err)
	}
	n, err := f.Write(firasB)

	if n != size {
		t.Errorf("Failed writing as io.Writer wrote %v bytes found %v bytes", n, size)
	}

	if err != nil {
		t.Fatal(err)
	}
}


func Test_Read(t *testing.T) {
	m := new(mockS3)
	SetService(m)

	m.On("GetObject").Return(&s3.GetObjectOutput{
		Body: new(ReadWriteCloser),
	}, nil)

	s3fs = New(Config{
		Id: "AKIAJBRFB4PEZIKTETJQ",
		Secret: "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",
		Region: "us-east-1",
		Bucket: "usman-gost",
	})

	f := s3fs.File("firas.jpg")

	//r, err := f.ReadShit()
	firasB, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	firas, err := os.Create("../../storage/firas_downloaded.jpg")
	if err != nil {
		t.Fatal(err)
	}

	firas.Write(firasB)
	firas.Close()
	//
}


func Test_GetString(t *testing.T) {
	_, err := s3fs.File("test.txt").ReadString()
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_GetSignedUrl(t *testing.T) {
	f := s3fs.File("test.txt")
	_, err := f.(*File).GetSignedUrl(time.Minute * 1)
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_Exist(t *testing.T)  {
	if ! s3fs.File("test.txt").Exist() {
		t.Errorf("File doesn't exist")
	}
}

func Test_Create_Dir(t *testing.T)  {
	if err := s3fs.Directory("dDir").Create(); err != nil {
		t.Errorf("Couldn't create dir")
	}
}

func Test_Exist_Dir(t *testing.T)  {
	if ! s3fs.Directory("dDir").Exist() {
		t.Errorf("Dir doesn't exist")
	}
}

func Test_Delete_Dir(t *testing.T)  {
	if err := s3fs.Directory("dDir").Delete(); err != nil {
		t.Errorf("Couldn't delete dir: %v", err)
	}
}

func Test_Dir_Not_Exist(t *testing.T)  {
	if s3fs.Directory("dDir").Exist() {
		t.Errorf("Dir exists")
	}
}

func Test_Stat(t *testing.T)  {
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
	err := s3fs.File("test.txt").Delete()
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
	_, err = s3fs.File("test.txt").ReadString()
	if err == nil {
		t.Errorf("File was not deleted in the bucket")
	}
}

func Test_NotExist(t *testing.T)  {
	if s3fs.File("test.txt").Exist() {
		t.Errorf("File does exist")
	}
}
