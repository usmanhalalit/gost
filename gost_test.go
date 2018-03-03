package gost

import (
	"testing"
	"fmt"
	"time"
)

func Test_New(t *testing.T) {
	_, err := New().PutString("test.txt", "abc")
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_GetString(t *testing.T) {
	_, err := New().GetString("test.txt")
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_GetSignedUrl(t *testing.T) {
	_, err := New().GetSignedUrl("test.txt", time.Minute * 1)
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}

func Test_Delete(t *testing.T) {
	err := New().Delete("test.txt")
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
	_, err = New().GetString("test.txt")
	if err == nil {
		t.Errorf("File was not deleted in the bucket")
	}
}
