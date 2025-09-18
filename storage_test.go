package main

import (
	"bytes"
	"io"
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

	key := "testdir"

	s := NewStore(config)

	data := []byte("some data to write to file")

	err := s.writeStream(key, bytes.NewReader(data))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	hasKey := s.Has(key)

	if !hasKey {
		t.Errorf("Expected To Have Key %s ", key)
	}

	r, err := s.Read(key)

	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("Expetced %s GOT %s", b, data)
	}

	s.Delete(key)
}
