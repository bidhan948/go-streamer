package main

import (
	"fmt"
	"log"

	"github.com/bidhan948/go-streamer/p2p"
)

func handlePeerConnection(p2p.Peer) error {
	fmt.Println("will implement something here later and from outside tcp transort :)")
	return nil
}

func main() {
	tcpConfig := p2p.TCPTransportConfig{
		ListenAddress: ":30090",
		HandshakeFunc: p2p.NOPHandShake,
		Decoder:       &p2p.DefaultDecoder{},
		OnPeerConnect: handlePeerConnection,
	}
	tr := p2p.NewTCPTransport(tcpConfig)

	go func() {
		for {
			msg := <-tr.Consume()
			log.Printf("Received message from %s: %+v\n", msg.From.String(), msg)
		}
	}()

	err := tr.ListenAndAccept()
	if err != nil {
		log.Fatal("SOMETHING WENT WRONG")
	}

	select {}
}
