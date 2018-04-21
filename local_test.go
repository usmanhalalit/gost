package gost

import (
	"github.com/usmanhalalit/gost/adapter/local"
	"log"
	"os"
	"testing"
	"time"
)

var ep, _ = os.Getwd()
var lfs = local.New(local.Config{
	BasePath: ep + "/storage",
})

func TestFiles(t *testing.T) {
	files, _ := lfs.Files()
	log.Println(files[1].GetPath())
}

func TestDirectories(t *testing.T) {
	dirs, _ := lfs.Directories()
	log.Println(dirs[0].GetPath())
}

func TestWrite(t *testing.T) {
	b := []byte("abc")
	n, err := lfs.File("test.txt").Write(b)
	if n != len(b) {
		t.Fatalf("Wrote %v bytes of %v bytes", n, len(b))
	}
	check(t, err)
}

func TestRead(t *testing.T) {
	b := make([]byte, 3)
	n, err := lfs.File("test.txt").Read(b)
	if n != len(b) {
		t.Fatalf("Read %v bytes of %v bytes", n, len(b))
	}
	check(t, err)
}

func TestStat(t *testing.T)  {
	info, err := lfs.File("test.txt").Stat()
	if err != nil {
		t.Errorf("Couldn't get stat: %v", err)
	}

	if info.Size != 3 {
		t.Errorf("Invalid file size expected %v got %v", 3, info.Size)
	}

	if info.LastModified.Day() != time.Now().Day() {
		t.Errorf("Invalid file time expected %v got %v", time.Now().Day(), info.LastModified.Day())
	}
}


func TestExist(t *testing.T) {
	if ! lfs.File("test.txt").Exist() {
		t.Fatalf("File does not exist")
	}
}


func TestExistDir(t *testing.T) {
	if ! lfs.Directory("aDir").Exist() {
		t.Fatalf("Dir does not exist")
	}
}

func TestDelete(t *testing.T) {
	check(t, lfs.File("test.txt").Delete())
}

func TestCreateDir(t *testing.T) {
	check(t, lfs.Directory("dDir").Create())
}

func TestStatDir(t *testing.T)  {
	info, err := lfs.Directory("dDir").Stat()
	if err != nil {
		t.Errorf("Couldn't get stat: %v", err)
	}

	if info.Size != 64 {
		t.Errorf("Invalid dir size expected %v got %v", 64, info.Size)
	}

	if info.LastModified.Day() != time.Now().Day() {
		t.Errorf("Invalid dir time expected %v got %v", time.Now().Day(), info.LastModified.Day())
	}
}

func TestDeleteDir(t *testing.T) {
	check(t, lfs.Directory("dDir").Delete())
}

func TestNotExist(t *testing.T) {
	if lfs.File("test.txt").Exist() {
		t.Fatalf("File exist")
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

