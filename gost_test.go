package gost

import (
	"fmt"
	"github.com/usmanhalalit/gost/adapter/s3"
	"io/ioutil"
	"os"
	"testing"
)

func Test_New(t *testing.T) {
	fs := New()
	//_, err := fs.File("test.txt").WriteString("abc")
	files, err := fs.Directory("aDir").Files()
	fmt.Println(files)
	fmt.Println(files[0].ReadString())
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}

	txt, err := fs.File("test.txt").ReadString()
	if err != nil {
		t.Errorf("Failed read: %v", err)
	}
	fmt.Printf(txt)
}

func Test_Write(t *testing.T) {
	fs := New()
	f := fs.File("firas.jpg")
	//n, err := fmt.Fprintf(f, "A formatted \na\na %v", "string")
	firas, err := os.Open("firas.jpg")
	if err != nil {
		t.Fatal(err)
	}
	fi, _ := firas.Stat()
	size := int(fi.Size())
	firasB := make([]byte, size)
	r, err := firas.Read(firasB)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Read %v bytes\n", r)
	fmt.Printf("Writing\n")
	n, err := f.Write(firasB)

	if n != size {
		t.Errorf("Failed writing as io.Writer wrote %v bytes found %v bytes", n, size)
	}

	if err != nil {
		t.Fatal(err)
	}
}


func Test_Read(t *testing.T) {
	fs := New()
	f := fs.File("firas.jpg").(*s3.S3file)

	//r, err := f.ReadShit()
	firasB, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	firas, err := os.Create("firas_downloaded.jpg")
	if err != nil {
		t.Fatal(err)
	}

	firas.Write(firasB)
	firas.Close()
	//
}


//func Test_GetString(t *testing.T) {
//	_, err := New().ReadString("test.txt")
//	if err != nil {
//		t.Errorf("Failed write: %v", err)
//	}
//}
//
//func Test_GetSignedUrl(t *testing.T) {
//	ad := New()
//	_, err := ad.(s3.S3filesystem).GetSignedUrl("test.txt", time.Minute * 1)
//	if err != nil {
//		t.Errorf("Failed write: %v", err)
//	}
//}
//
//
//func Test_Exist(t *testing.T)  {
//	if ! New().Exist("test.txt") {
//		t.Errorf("File doesn't exist")
//	}
//}
//
//func Test_Info(t *testing.T)  {
//	info, err := New().Info("test.txt")
//	if err != nil {
//		t.Errorf("Couldn't get info: %v", err)
//	}
//
//	if info.Size != 3 {
//		t.Errorf("Invalid file size expected %v got %v", 3, info.Size)
//	}
//
//	if info.LastModified.Day() != time.Now().Day() {
//		t.Errorf("Invalid file size expected %v got %v", time.Now().Day(), info.LastModified.Day())
//	}
//}
//
//func Test_Delete(t *testing.T) {
//	err := New().Delete("test.txt")
//	if err != nil {
//		t.Errorf("Failed write: %v", err)
//	}
//	_, err = New().ReadString("test.txt")
//	if err == nil {
//		t.Errorf("File was not deleted in the bucket")
//	}
//}
//
//func Test_NotExist(t *testing.T)  {
//	if New().Exist("test.txt") {
//		t.Errorf("File does exist")
//	}
//}
