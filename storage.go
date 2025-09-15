package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type PathTransformFunc func(string) string

type StoreConfig struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreConfig
}

func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen+1)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

func NewStore(config StoreConfig) *Store {
	return &Store{
		StoreConfig: config,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Println("error creating directory", err)
		return err
	}

	filename := "some_random"
	fullPathAndFilename := path + "/" + filename
	f, err := os.Create(fullPathAndFilename)

	if err != nil {
		fmt.Println("error opening file", err)
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		fmt.Println("error writing to file", err)
		return err
	}
	fmt.Println("wrote", n, "bytes to", fullPathAndFilename)
	return nil
}
