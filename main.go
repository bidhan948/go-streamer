package main

import (
	"log"

	"github.com/bidhan948/go-streamer/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":30090")
	err := tr.ListenAndAccept()
	if err != nil {
		log.Fatal("SOMETHING WENT WRONG")
	}

	select {}
}
