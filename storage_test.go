package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	keyName := "SomeRandomKey"
	pathKey := CASPathTransformFunc(keyName)
	expectedPath := "dc2a1/164e7/b8e5a/d09fb/8e12c/4040e/584a3/e5867/"
	expectedOrginalPathName := "dc2a1164e7b8e5ad09fb8e12c4040e584a3e5867"

	if pathKey.Pathname != expectedPath {
		t.Fatalf("expected %s, got %s", expectedPath, pathKey.Pathname)
	}

	if pathKey.FileName != expectedOrginalPathName {
		t.Fatalf("expected %s, got %s", expectedOrginalPathName, pathKey.FileName)
	}
}

func TestStore(t *testing.T) {
	config := StoreConfig{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(config)

	data := bytes.NewReader([]byte("some data to write to file"))

	err := s.writeStream("testdir", data)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
