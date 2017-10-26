package packet

import (
	"github.com/bamiaux/iobit"
	"fmt"
	"errors"
)

const PES_PACKET_START_CODE_PREFIX = 0x000001

type PESFlags  struct {iobit.Reader}

type PES struct {
	iobit.Reader
	b []byte
}

func IsPES(b []byte) bool {
	if len(b) == 0 { return false }
	r := iobit.NewReader(b)
	start := r.Uint32(24)
	return start == PES_PACKET_START_CODE_PREFIX
}

func NewPES(b []byte) (p *PES, err error) {
	p = &PES{iobit.NewReader(b), b}
	start := p.Uint32(24)
	if start != PES_PACKET_START_CODE_PREFIX {
		return nil, errors.New(fmt.Sprintf("Invalid PES Packet - Incorrect packet start code prefix. Expected %0x but found %0x.", PES_PACKET_START_CODE_PREFIX, start))
	}
	return
}

func (p *PES) Flags() *PESFlags {
	return &PESFlags{p.Reader}
}

func (p *PESFlags) PTSDTS() uint8 {
	p.Reset()
	p.Skip(56)
	return p.Uint8(2)
}

func (p *PESFlags) PTS() bool {
	p.Reset()
	p.Skip(56)
	return p.Uint8(2) >= 2
}

func (p *PESFlags) DTS() bool {
	p.Reset()
	p.Skip(56)
	return p.Uint8(2) >= 3
}

func (p *PES) PTS() (pts uint64, e error) {
	f := PESFlags {p.Reader}

	if !f.PTS() {
		return 0, errors.New("Missing pts_flag.")
	}

	p.Reset()


	p.Skip(70)

	var v uint64 = p.Uint64(36) // this is a 64bit integer, lowest 36 bits contain a timestamp with markers
	pts = 0
	pts |= (v >> 3) & (0x0007 << 30) // top 3 bits, shifted left by 3, other bits zeroed out
	pts |= (v >> 2) & (0x7fff << 15) // middle 15 bits
	pts |= (v >> 1) & (0x7fff <<  0) // bottom 15 bits

	pts = pts

	return
}

func (p *PES) DTS() (dts uint64, e error) {
	f := PESFlags {p.Reader}

	if !f.DTS() {
		return 0, errors.New("Missing dts_flag.")
	}

	p.Reset()

	p.Skip(70)


	if f.PTS() {
		p.Skip(36)
	}

	var v uint64 = p.Uint64(36) // this is a 64bit integer, lowest 36 bits contain a timestamp with markers
	dts = 0
	dts |= (v >> 3) & (0x0007 << 30) // top 3 bits, shifted left by 3, other bits zeroed out
	dts |= (v >> 2) & (0x7fff << 15) // middle 15 bits
	dts |= (v >> 1) & (0x7fff <<  0) // bottom 15 bits

	dts = dts

	return
}