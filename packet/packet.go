package packet

import (
	"errors"
	"github.com/bamiaux/iobit"
	"fmt"
)

const SYNC_BYTE_MAGIC = 0x47

type Packet struct {
	TEI uint8
	PUSI uint8
	TransportPriority uint8
	PID uint16
	TSC uint8
	AdaptionFieldControl uint8
	ContinuityCounter uint8
	AdaptionField *AdaptionFieldPayload
	Payload *[]byte
}

type AdaptionFieldPayload struct {
	Discontinuity bool
	RandomAccess bool
	ElementryStreamPriority bool
	Flags AdaptionFieldFlags
	OptionalFields AdaptionFieldOptionalFields
}

type AdaptionFieldFlags struct {
	PCR bool
	OPCR bool
	SplicingPoint bool
	transportPrivateData bool
	AdaptionFieldExtension bool
}

type AdaptionFieldOptionalFields struct {
	PCR *uint64
	OPCR *uint64
	SpliceCountdown *uint8
	TransportPrivateData *[]byte
	Extension *AdaptionFieldExtension
}

type AdaptionFieldExtension struct {
	LTWFlag bool
	PiecewiseRateFlag bool
	SeamlessSpliceFlag bool
	LTW *AFLegalTimeWindow
	PiecewiseRate *uint32
	SeamlessSplice *AFSeamlessSplice
}

type AFLegalTimeWindow struct {
	valid bool
	offset uint16
}

type AFSeamlessSplice struct {
	SpliceType uint8
	DTSNextAUs []AFSeamlessSpliceDTSNextAU
}

type AFSeamlessSpliceDTSNextAU struct {
	NextAU uint16
	Marker bool
}


func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		bnoden = str[0:num]
	}
	return bnoden
}


func DecodeTimestamp(ts uint64, clk uint64) (str string) {
	//var fa int64 = 27000000;
	var h uint64 = (ts/(clk*60*60));
	var m uint64 = (ts/(clk*60))-(h*60);
	var s uint64 = (ts/clk)-(h*3600)-(m*60);
	var u uint64 = ts-(h*clk*60*60)-(m*clk*60)-(s*clk);
	return fmt.Sprintf(" %02d:%02d:%02d.%s", h, m, s, truncateString(fmt.Sprintf("%05d", u), 5))
}


func NewPacket(b []byte) ( *Packet, error) {
	r := iobit.NewReader(b)

	if r.Byte() != SYNC_BYTE_MAGIC {

		return nil, errors.New("Invalid TS Packet - Incorrect sync byte.")
	}

	tei := r.Uint8(1)
	pusi := r.Uint8(1)
	tp := r.Uint8(1)

	pid := r.Uint16(13)

	tsc := r.Uint8(2)
	afc := r.Uint8(2)

	cc := r.Uint8(4)


	var afp *AdaptionFieldPayload = nil
	var afp_len byte = 0

	if afc == 2 || afc == 3 { // af cntl 2||3 && af len > 0
		afp_len = r.Byte()
		if afp_len > 0 {
			df := r.Bit()
			raf := r.Bit()
			espf := r.Bit()
			pcrf := r.Bit()
			opcrf := r.Bit()
			spf := r.Bit()
			tpdf := r.Bit()
			afef := r.Bit()

			aff := AdaptionFieldFlags{
				pcrf,
				opcrf,
				spf,
				tpdf,
				afef,
			}

			var pcr *uint64 = nil
			var opcr *uint64 = nil
			var sc *uint8 = nil
			var tpd *[]byte = nil
			var ext *AdaptionFieldExtension = nil

			if pcrf {
				base := r.Uint64(33) // PCR base is 33-bits
				r.Skip(6)           // 6-bits are reserved
				ext := r.Uint64(9)  // PCR extension is 9-bits
				pcr_val := base*300 + ext
				pcr = &pcr_val
			}

			if opcrf {
				base := r.Uint64(33) // PCR base is 33-bits
				r.Skip(6)           // 6-bits are reserved
				ext := r.Uint64(9)  // PCR extension is 9-bits
				opcr_val := base*300 + ext
				opcr = &opcr_val
			}

			if spf {
				sp_val := r.Uint8(8)
				sc = &sp_val
			}

			if tpdf {
				tpdLength := r.Uint8(8)
				tpd_val := make([]byte, tpdLength*8, tpdLength*8)
				for i := 0; i < int(tpdLength*8); i++ {
					tpd_val[i] = r.Byte()
				}
				tpd = &tpd_val
			}

			if afef {
				ltwf := false
				prf := false
				ssf := false

				var ltw *AFLegalTimeWindow = nil

				var pr *uint32 = nil
				var ss *AFSeamlessSplice = nil

				if ltwf {
					ltw_v := false
					ltw_o := uint16(0)
					ltw = &AFLegalTimeWindow{
						ltw_v,
						ltw_o,
					}
				}

				if prf {
					pr_val := uint32(0)
					pr = &pr_val
				}

				if ssf {
					ss_st := uint8(0)
					dtsnaus := make([]AFSeamlessSpliceDTSNextAU, 3, 3)

					dtsnaus[0].NextAU = 0
					dtsnaus[0].Marker = false

					dtsnaus[1].NextAU = 0
					dtsnaus[1].Marker = false

					dtsnaus[2].NextAU = 0
					dtsnaus[2].Marker = false

					ss = &AFSeamlessSplice{
						ss_st,
						dtsnaus,
					}
				}

				ext = &AdaptionFieldExtension{
					ltwf,
					prf,
					ssf,
					ltw,
					pr,
					ss,
				}
			}

			afop := AdaptionFieldOptionalFields{
				pcr,
				opcr,
				sc,
				tpd,
				ext,
			}

			afp = &AdaptionFieldPayload{
				df,
				raf,
				espf,
				aff,
				afop,
			}
		}
	} else {
		afc = 0 // for the odd case where afc is set to 2 or 3 but the af len is 0
	}

	payloadStart := 4  // fixed length part of packet header

	if afc == 2 || afc == 3 {
		payloadStart += 1 // afp_len doesn't inc itself
	}

	payloadStart += int(afp_len) //dynamic length part of packet header

	//fmt.Println("payloadStart", payloadStart*8, r.At())

	r.Skip(uint(payloadStart*8)-r.At())

	//pl := r.Bytes(188-payloadStart)
	var pl []byte = nil
	return &Packet{
		tei,
		pusi,
		tp,
		pid,
		tsc,
		afc,
		cc,
		afp,
		&pl,
	}, nil
}

/*func ( p *Packet) String() string {
	return "t"
}*/