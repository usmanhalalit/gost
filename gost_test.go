package gost

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	fs := New()
	//_, err := fs.File("test.txt").PutString("abc")
	files, err := fs.Files()
	fmt.Println(files)
	fmt.Println(files[0].GetString())
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}

	txt, err := fs.File("test.txt").GetString()
	if err != nil {
		t.Errorf("Failed read: %v", err)
	}
	fmt.Printf(txt)
}

//
//func Test_NewIdea(t *testing.T) {
//	disc := New()
//	file := disc.File("test.txt")
//	_, err := file.PutString("abc")
//	if err != nil {
//		t.Errorf("Failed write: %v", err)
//	}
//}
//
//func Test_GetString(t *testing.T) {
//	_, err := New().GetString("test.txt")
//	if err != nil {
//		t.Errorf("Failed write: %v", err)
//	}
//}
//
//func Test_GetSignedUrl(t *testing.T) {
//	ad := New()
//	_, err := ad.(s3.S3adapter).GetSignedUrl("test.txt", time.Minute * 1)
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
//	_, err = New().GetString("test.txt")
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
