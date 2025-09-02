package p2p

type handshakeFunc func(any) error

func NOPHandShake(any) error { return nil }
