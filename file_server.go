package main

import "github.com/bidhan948/go-streamer/p2p"

type FileServerConfig struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerConfig
	Store *Store
}

func NewFileServer(config FileServerConfig) *FileServer {
	storeConfig := StoreConfig{
		Root:              config.StorageRoot,
		PathTransformFunc: config.PathTransformFunc,
	}

	return &FileServer{
		FileServerConfig: config,
		Store:            NewStore(storeConfig),
	}
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}
