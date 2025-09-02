package p2p

import "io"

type Decoder interface {
	decode(io.Reader, any) error
}
