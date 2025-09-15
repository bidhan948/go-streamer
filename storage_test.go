package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	keyName := "SomeRandomKey"
	pathName := CASPathTransformFunc(keyName)
	expectedPath := "dc2a1/164e7/b8e5a/d09fb/8e12c/4040e/584a3/e5867/"

	if pathName != expectedPath {
		t.Fatalf("expected %s, got %s", expectedPath, pathName)
	}
}

func TestStore(t *testing.T) {
	config := StoreConfig{
		PathTransformFunc: DefaultPathTransformFunc,
	}

	s := NewStore(config)

	data := bytes.NewReader([]byte("some data to write to file"))

	err := s.writeStream("testdir", data)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
