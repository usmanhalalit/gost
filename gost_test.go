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
	url, err := New().GetSignedUrl("test.txt", time.Minute * 1)
	fmt.Println(url)
	if err != nil {
		t.Errorf("Failed write: %v", err)
	}
}
