package main

import (
	"log"

	"github.com/bidhan948/go-streamer/p2p"
)

func main() {
	tcpConfig := p2p.TCPTransportConfig{ListenAddress: ":30090", HandshakeFunc: p2p.NOPHandShake, Decoder: p2p.GOBDecoder{}}
	tr := p2p.NewTCPTransport(tcpConfig)

	err := tr.ListenAndAccept()
	if err != nil {
		log.Fatal("SOMETHING WENT WRONG")
	}

	select {}
}
