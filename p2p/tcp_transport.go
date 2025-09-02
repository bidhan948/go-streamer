package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshakeFunc handshakeFunc
	decoder       Decoder

	mu    sync.RWMutex
	peers map[string]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
		handshakeFunc: NOPHandShake,
	}
}

func NewTCPPeer(conn net.Conn, outbond bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbond,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
		}
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	err := t.handshakeFunc(conn)
	if err != nil {
		log.Fatal("CANNOT SHAKE HANDS : ", err)
	}
	fmt.Printf("NEW INCMING CONNECTION : %+v", peer)
}
