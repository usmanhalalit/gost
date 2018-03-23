package gost

import (
	"github.com/usmanhalalit/gost/adapter/s3"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var s3fs = New()

func Test_New(t *testing.T) {
	files, err := s3fs.Directory("aDir").Files()

	if len(files) < 1 {
		t.Errorf("Failed listing found %v files", len(files))
	}

	if err != nil {
		t.Errorf("Failed listing: %v", err)
	}

	err = s3fs.File("test.txt").WriteString("abc")
	if err != nil {
		t.Errorf("Failed write: %v", err)
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
	_, err = New().File("test.txt").ReadString()
	if err == nil {
		t.Errorf("File was not deleted in the bucket")
	}
}

func Test_NotExist(t *testing.T)  {
	if s3fs.File("test.txt").Exist() {
		t.Errorf("File does exist")
	}
}
