package impulse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type SampleFlag uint8

const (
	SampleAssociatedWithHeader SampleFlag = 1 << iota
	Quality16Bit
	StereoSample
	Compressed
	Loop
	SustainLoop
	PingPongLoop
	PingPongSustainLoop
)

type VibratoWaveform uint8

const (
	SineWave VibratoWaveform = iota
	RampDown
	SquareWave
	Random
)

type rawSample struct {
	MagicString        [4]byte
	DOSFilename        [12]byte
	_                  byte
	GvL, Flg, Vol      uint8
	SampleName         [26]byte
	Cvt, DfP           uint8
	Length             uint32
	LoopBegin          uint32
	LoopEnd            uint32
	C5Speed            uint32
	SusLoopBegin       uint32
	SusLoopEnd         uint32
	SamplePointer      uint32
	ViS, ViD, ViR, ViT uint8
}

// Sample is an Impulse Tracker sample.
type Sample struct {
	Filename         string // max 12 bytes
	GlobalVolume     uint8  // range 0->64
	Flags            SampleFlag
	DefaultVolume    uint8  // range 0->64
	Name             string // max 26 bytes
	Signed           bool
	DefaultPan       uint8 // range 0->64
	DefaultPanOn     bool
	Length           uint32
	LoopBegin        uint32
	LoopEnd          uint32
	Speed            uint32 // range 0->9999999
	SustainLoopBegin uint32
	SustainLoopEnd   uint32
	VibratoSpeed     uint8 // range 0->64
	VibratoDepth     uint8 // range 0->64
	VibratoWaveform  VibratoWaveform
	VibratoRate      uint8
	Data             []byte // PCM audio data
}

func sampleFromRaw(raw *rawSample, r io.ReadSeeker) (*Sample, error) {
	s := Sample{
		Filename:         string(bytes.Trim(raw.DOSFilename[:], "\x00")),
		GlobalVolume:     raw.GvL,
		Flags:            SampleFlag(raw.Flg),
		DefaultVolume:    raw.Vol,
		Name:             string(bytes.Trim(raw.SampleName[:], "\x00")),
		Signed:           raw.Cvt&0x01 != 0,
		DefaultPan:       raw.DfP & 0x7f,
		DefaultPanOn:     raw.DfP&0x80 != 0,
		Length:           raw.Length,
		LoopBegin:        raw.LoopBegin,
		LoopEnd:          raw.LoopEnd,
		Speed:            raw.C5Speed,
		SustainLoopBegin: raw.SusLoopBegin,
		SustainLoopEnd:   raw.SusLoopEnd,
		VibratoSpeed:     raw.ViS,
		VibratoDepth:     raw.ViD,
		VibratoRate:      raw.ViR,
		VibratoWaveform:  VibratoWaveform(raw.ViT),
	}

	if s.Flags&Quality16Bit != 0 {
		s.Data = make([]byte, s.Length*2)
	} else {
		s.Data = make([]byte, s.Length)
	}

	if _, err := r.Seek(int64(raw.SamplePointer), 0); err != nil {
		return nil, err
	}
	if _, err := r.Read(s.Data); err != nil {
		return nil, err
	}

	return &s, nil
}

// ReadSample reads a Sample in ITS format from r.
func ReadSample(r io.ReadSeeker) (*Sample, error) {
	raw := new(rawSample)
	if err := binary.Read(r, binary.LittleEndian, raw); err != nil {
		return nil, err
	}
	if string(raw.MagicString[:]) != "IMPS" {
		return nil, errors.New("data is not Impulse Tracker sample")
	}
	return sampleFromRaw(raw, r)
}

// Write writes the Sample to w in ITS format.
func (s *Sample) Write(w io.Writer) error {
	raw := &rawSample{
		MagicString:   [4]byte{'I', 'M', 'P', 'S'},
		GvL:           s.GlobalVolume,
		Flg:           uint8(s.Flags),
		Vol:           s.DefaultVolume,
		DfP:           s.DefaultPan,
		Length:        s.Length,
		LoopBegin:     s.LoopBegin,
		LoopEnd:       s.LoopEnd,
		C5Speed:       s.Speed,
		SusLoopBegin:  s.SustainLoopBegin,
		SusLoopEnd:    s.SustainLoopEnd,
		SamplePointer: 0x0050,
		ViS:           s.VibratoSpeed,
		ViD:           s.VibratoDepth,
		ViR:           s.VibratoRate,
		ViT:           uint8(s.VibratoWaveform),
	}
	for i := range raw.DOSFilename {
		if i < len(s.Filename) {
			raw.DOSFilename[i] = s.Filename[i]
		}
	}
	for i := range raw.SampleName {
		if i < len(s.Name) {
			raw.SampleName[i] = s.Name[i]
		}
	}
	if s.Signed {
		raw.Cvt = 0x01
	}
	if s.DefaultPanOn {
		raw.DfP |= 0x80
	}

	if err := binary.Write(w, binary.LittleEndian, raw); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, s.Data); err != nil {
		return err
	}

	return nil
}
