package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultDirName = "ggnetwork"

type PathTransformFunc func(string) PathKey

type StoreConfig struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreConfig
}

type PathKey struct {
	Pathname string
	FileName string
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		FileName: key,
	}
}

func NewStore(config StoreConfig) *Store {
	if config.PathTransformFunc == nil {
		config.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(config.Root) == 0 {
		config.Root = defaultDirName
	}
	return &Store{
		StoreConfig: config,
	}
}

func (p PathKey) FirstPathName() string {
	path := strings.Split(p.Pathname, "/")

	if len(path) == 0 {
		return ""
	}

	return path[0]
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	filePathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FilePath())
	_, err := os.Stat(filePathWithRoot)

	return err != fs.ErrNotExist
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted %s from disk", pathKey.FileName)
	}()
	firstPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	return os.RemoveAll(firstPathWithRoot)
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen+1)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

func (p PathKey) FilePath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.FileName)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, path.Pathname)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		fmt.Println("error creating directory", err)
		return err
	}

	fullPathAndFilename := path.FilePath()
	fullPathAndFilenameWithRoot := fmt.Sprintf("%s/%s", s.Root, fullPathAndFilename)

	f, err := os.Create(fullPathAndFilenameWithRoot)

	if err != nil {
		fmt.Println("error opening file", err)
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		fmt.Println("error writing to file", err)
		return err
	}

	fmt.Println("wrote", n, "bytes to", fullPathAndFilenameWithRoot)
	return nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FilePath())
	return os.Open(fullPathWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := new(bytes.Buffer)

	_, err = io.Copy(buff, f)

	if err != nil {
		return nil, err
	}

	return buff, nil
}
