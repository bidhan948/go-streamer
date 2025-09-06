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

type TCPTransportConfig struct {
	ListenAddress string
	HandshakeFunc handshakeFunc
	Decoder       Decoder
}
type TCPTransport struct {
	TCPTransportConfig
	listener net.Listener
	mu       sync.RWMutex
	peers    map[string]Peer
}

func NewTCPTransport(config TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		TCPTransportConfig: config,
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
	t.listener, err = net.Listen("tcp", t.ListenAddress)

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
		fmt.Printf("NEW INCOMING CONNECTION : %+v /n", conn)
		go t.handleConnection(conn)

	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	err := t.HandshakeFunc(peer)
	if err != nil {
		log.Fatal("CANNOT SHAKE HANDS : ", err)
		conn.Close()
		return
	}

	buff := make([]byte, 2000)
	for {
		n, err := conn.Read(buff)

		if err != nil {
			log.Fatal("Buffer Err ", err)
		}

		fmt.Printf("MESSAGE : %+v\n", buff[:n])
	}
}
