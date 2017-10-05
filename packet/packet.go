package packet

import (
	"github.com/bamiaux/iobit"
	"errors"
	"fmt"
)
//import "github.com/bamiaux/iobit"

const SYNC_BYTE_MAGIC = 0x47

type Packet struct {iobit.Reader}

func NewPacket(b []byte) (p *Packet, err error) {
	p = &Packet{iobit.NewReader(b)}
	if p.Byte() != SYNC_BYTE_MAGIC {
		return nil, errors.New("Invalid TS Packet - Incorrect sync byte.")
	}
	return
}

func (p *Packet) TEI() bool {
	p.Reset()
	p.Skip(8)
	return p.Bit()
}

func (p *Packet) PUSI() bool {
	p.Reset()
	p.Skip(9)
	return p.Bit()
}

func (p *Packet) TransportPriority() bool {
	p.Reset()
	p.Skip(10)
	return p.Bit()
}

func (p *Packet) PID() uint16 {
	p.Reset()
	p.Skip(11)
	return p.Uint16(13)
}

func (p *Packet) TSC() uint8 {
	p.Reset()
	p.Skip(24)
	return p.Uint8(2)
}

func (p *Packet) AdaptionFieldControl() uint8 {
	p.Reset()
	p.Skip(26)
	return p.Uint8(2)
}

func (p *Packet) ContinuityCounter() uint8 {
	p.Reset()
	p.Skip(28)
	return p.Uint8(4)
}

func (p Packet) AdaptionField() *AdaptionField {
	af := AdaptionField{p.Reader}
	return &af
}

func DecodeTimestamp(ts uint64, p uint64, fa uint64) (str string) {
	var (
		h uint64
		m uint64
		s uint64
		u uint64
	)
	ts /= uint64(p) // Convert to milliseconds
	h = (ts / (fa * 60 * 60))
	m = (ts / (fa * 60)) - (h * 60)
	s = (ts / fa) - (h * 3600) - (m * 60)
	u = ts - (h * fa * 60 * 60) - (m * fa * 60) - (s * fa)
	return fmt.Sprintf("%02dh%02dm%02ds%dµs", h, m, s, u)
}
