package p2p

import (
	"fmt"
	"log"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

type TCPTransportConfig struct {
	ListenAddress string
	HandshakeFunc handshakeFunc
	Decoder       *DefaultDecoder
	OnPeerConnect func(Peer) error
}

type TCPTransport struct {
	TCPTransportConfig
	listener net.Listener
	rpch     chan *RPC
}

func NewTCPTransport(config TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		TCPTransportConfig: config,
		rpch:               make(chan *RPC),
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) Consume() <-chan *RPC {
	return t.rpch
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
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
	var err error

	defer func() {
		fmt.Printf("CLOSING CONNECTION : %+v /n", conn)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	err = t.HandshakeFunc(peer)
	if err != nil {
		log.Fatal("CANNOT SHAKE HANDS : ", err)
		return
	}

	if t.OnPeerConnect != nil {
		if err := t.OnPeerConnect(peer); err != nil {
			log.Println("OnPeerConnect error", err)
			return
		}
	}
	rpc := &RPC{}

	for {
		err := t.Decoder.Decode(conn, rpc)
		if err != nil {
			log.Println("Error decoding message", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpch <- rpc
	}
}
