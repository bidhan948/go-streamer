package p2p

type handshakeFunc func(Peer) error

func NOPHandShake(Peer) error { return nil }
