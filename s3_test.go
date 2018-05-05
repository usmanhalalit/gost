package gost

import (
	"github.com/usmanhalalit/gost/adapter/s3"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

var s3fs = s3.NewS3Adapter(s3.S3config{
	Id: "AKIAJBRFB4PEZIKTETJQ",
	Secret: "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",
	Region: "us-east-1",
	Bucket: "usman-gost",
})

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
	//n, err := fmt.Fprintf(f, "A formatted \na\na %v", "string")
	firas, err := os.Open("storage/firas.jpg")
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
	f := s3fs.File("firas.jpg").(*s3.S3file)

	//r, err := f.ReadShit()
	firasB, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	firas, err := os.Create("storage/firas_downloaded.jpg")
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
	_, err := f.(*s3.S3file).GetSignedUrl(time.Minute * 1)
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_Exist(t *testing.T)  {
	if ! s3fs.File("test.txt").Exist() {
		t.Errorf("File doesn't exist")
	}
}

func Test_Exist_Dir(t *testing.T)  {
	if ! s3fs.Directory("aDir").Exist() {
		t.Errorf("Dir doesn't exist")
	}
}

func Test_Dir_Not_Exist(t *testing.T)  {
	if s3fs.Directory("xDir").Exist() {
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
