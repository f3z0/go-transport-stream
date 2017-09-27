package ts

import "github.com/f3z0/ts/packet"

type Reader struct {
}

func (r *Reader) Read(p []packet.Packet) (n int, err error) {
	return 0, nil
}

