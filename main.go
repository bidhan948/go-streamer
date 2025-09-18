package main

import (
	"log"

	"github.com/bidhan948/go-streamer/p2p"
)

func main() {

	tcpTransportConfig := p2p.TCPTransportConfig{
		ListenAddress: ":30090",
		HandshakeFunc: p2p.NOPHandShake,
		Decoder:       &p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportConfig)

	fileServerConfig := FileServerConfig{
		StorageRoot:       "30090_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fs := NewFileServer(fileServerConfig)

	if err := fs.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
