package packet

import (
	"github.com/bamiaux/iobit"
	"errors"
)

type AdaptionField struct {iobit.Reader}
type AdaptionFieldFlags  struct {iobit.Reader}
type AdaptionFieldOptionalFields  struct {iobit.Reader}
type AdaptionFieldExtension  struct {iobit.Reader}
type AFLegalTimeWindow  struct {iobit.Reader}
type AFSeamlessSplice  struct {iobit.Reader}
type AFSeamlessSpliceDTSNextAU  struct {iobit.Reader}

func (p *AdaptionField) Discontinuity() bool {
	p.Reset()
	p.Skip(40)
	return p.Bit()
}

func (p *AdaptionField) RandomAccess() bool {
	p.Reset()
	p.Skip(41)
	return p.Bit()
}

func (p *AdaptionField) ElementryStreamPriority() bool {
	p.Reset()
	p.Skip(42)
	return p.Bit()
}

func (p *AdaptionField) Flags() *AdaptionFieldFlags {
	return &AdaptionFieldFlags{p.Reader}
}

func (p *AdaptionField) OptionalFields() *AdaptionFieldOptionalFields {
	return &AdaptionFieldOptionalFields{p.Reader}
}

func (p *AdaptionFieldFlags) PCR() bool {
	p.Reset()
	p.Skip(43)
	return p.Bit()
}

func (p *AdaptionFieldFlags) OPCR() bool {
	p.Reset()
	p.Skip(44)
	return p.Bit()
}

func (p *AdaptionFieldFlags) SplicingPoint() bool {
	p.Reset()
	p.Skip(45)
	return p.Bit()
}

func (p *AdaptionFieldFlags) TransportPrivateData() bool {
	p.Reset()
	p.Skip(46)
	return p.Bit()
}

func (p *AdaptionFieldFlags) AdaptionFieldExtension() bool {
	p.Reset()
	p.Skip(47)
	return p.Bit()
}

func (p *AdaptionFieldOptionalFields) PCR() (pcr uint64, e error) {
	f := AdaptionFieldFlags {p.Reader}

	if !f.PCR() {
		return 0, errors.New("Missing pcr_flag.")
	}

	p.Reset()

	p.Skip(48)
	base := p.Uint64(33) // PCR base is 33-bits
	p.Skip(6)           // 6-bits are reserved
	ext := p.Uint64(9)  // PCR extension is 9-bits
	pcr = base*300 + ext
	return
}

func (p *AdaptionFieldOptionalFields) OPCR() (pcr uint64, e error) {
	f := AdaptionFieldFlags {p.Reader}

	if !f.OPCR() {
		return 0, errors.New("Missing opcr_flag.")
	}

	p.Reset()

	p.Skip(48)

	if f.PCR() {
		p.Skip(48)
	}

	base := p.Uint64(33) // PCR base is 33-bits
	p.Skip(6)           // 6-bits are reserved
	ext := p.Uint64(9)  // PCR extension is 9-bits
	pcr = base*300 + ext
	return
}

func (p *AdaptionFieldOptionalFields) SplicingPoint() (sp uint8, e error) {
	f := AdaptionFieldFlags {p.Reader}

	if !f.SplicingPoint() {
		return 0, errors.New("Missing splicing_point_flag.")
	}

	p.Reset()

	p.Skip(48)

	if f.PCR() {
		p.Skip(48)
	}

	if f.OPCR() {
		p.Skip(48)
	}

	sp = p.Uint8(8)
	return
}

func (p *AdaptionFieldOptionalFields) TransportPrivateData() (tpd []byte, e error) {
	f := AdaptionFieldFlags {p.Reader}

	if !f.TransportPrivateData() {
		return nil, errors.New("Missing transport_private_data_flag.")
	}

	p.Reset()

	p.Skip(48)

	if f.PCR() {
		p.Skip(48)
	}

	if f.OPCR() {
		p.Skip(48)
	}

	if f.SplicingPoint() {
		p.Skip(8)
	}

	tpdLength := p.Uint8(8)
	tpd = make([]byte, tpdLength*8, tpdLength*8)
	for i := 0; i < int(tpdLength*8); i++ {
		tpd[i] = p.Byte()
	}
	return
}

func (p *AdaptionFieldOptionalFields) Extension() *AdaptionFieldExtension {
	return nil
}

func (p *AdaptionFieldExtension) LTWFlag() bool {
	return false
}

func (p *AdaptionFieldExtension) PiecewiseRateFlag() bool {
	return false
}

func (p *AdaptionFieldExtension) SeamlessSpliceFlag() bool {
	return false
}

func (p *AdaptionFieldExtension) LTW() *AFLegalTimeWindow {
	return nil
}

func (p *AdaptionFieldExtension) PiecewiseRate() uint32 {
	return 0
}

func (p *AdaptionFieldExtension) SeamlessSplice() *AFSeamlessSplice {
	return nil
}

func (p *AFLegalTimeWindow) Valid() bool {
	return false
}

func (p *AFLegalTimeWindow) Offset() uint16 {
	return 0
}

func (p *AFSeamlessSplice) SpliceType() uint8 {
	return 0
}

func (p *AFSeamlessSplice) DTSNextAUs() []AFSeamlessSpliceDTSNextAU {
		return nil
}

func (p *AFSeamlessSpliceDTSNextAU) NextAU() uint16 {
	return 0
}

func (p *AFSeamlessSpliceDTSNextAU) Marker() bool {
	return false
}
