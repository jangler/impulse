package impulse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type rawModule struct {
	MagicString   [4]byte
	SongName      [26]byte
	PHilight      [2]byte
	OrdNum        uint16
	InsNum        uint16
	SmpNum        uint16
	PatNum        uint16
	Cwtv          [2]byte
	Cmwt          [2]byte
	Flags         uint16
	Special       uint16
	GV, MV        uint8
	IS, IT        uint8
	Sep, PWD      uint8
	MsgLgth       uint16
	MessageOffset uint32
	_             uint32
	ChnlPan       [64]byte
	ChnlVol       [64]byte
}

// Module is an Impulse Tracker module.
type Module struct {
	SongName        string // max 26 bytes
	GlobalVolume    uint8  // range 0->128
	MixingVolume    uint8  // range 0->128
	InitialSpeed    uint8
	InitialTempo    uint8
	Separation      uint8 // range 0->128
	PitchWheelDepth uint8
	ChannelPanning  []uint8 // range 0->64
	ChannelVolume   []uint8 // range 0->64
	OrderList       []uint8 // range 0->199, 254, 255
	Samples         []*Sample
}

func moduleFromRaw(raw *rawModule, r io.ReadSeeker) (*Module, error) {
	m := &Module{
		SongName:        string(bytes.Trim(raw.SongName[:], "\x00")),
		GlobalVolume:    raw.GV,
		MixingVolume:    raw.MV,
		InitialSpeed:    raw.IS,
		InitialTempo:    raw.IT,
		Separation:      raw.Sep,
		PitchWheelDepth: raw.PWD,
		ChannelPanning:  make([]uint8, 64),
		ChannelVolume:   make([]uint8, 64),
		OrderList:       make([]uint8, raw.OrdNum),
		Samples:         make([]*Sample, raw.SmpNum),
	}

	for i := range m.ChannelPanning {
		m.ChannelPanning[i] = raw.ChnlPan[i]
	}
	for i := range m.ChannelVolume {
		m.ChannelVolume[i] = raw.ChnlVol[i]
	}
	if _, err := r.Read(m.OrderList); err != nil {
		return nil, err
	}
	for i := range m.Samples {
		var err error
		ptrOffset := 0xc0 + int64(raw.OrdNum+raw.InsNum*4) + int64(i*4)
		if _, err = r.Seek(ptrOffset, 0); err != nil {
			return nil, err
		}
		var smpOffset uint32
		binary.Read(r, binary.LittleEndian, &smpOffset)
		if _, err = r.Seek(int64(smpOffset), 0); err != nil {
			return nil, err
		}
		if m.Samples[i], err = ReadSample(r); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// ReadModule reads a Module in IT format from r.
func ReadModule(r io.ReadSeeker) (*Module, error) {
	raw := new(rawModule)
	if err := binary.Read(r, binary.LittleEndian, raw); err != nil {
		return nil, err
	}
	if string(raw.MagicString[:]) != "IMPM" {
		return nil, errors.New("data is not Impulse Tracker module")
	}
	return moduleFromRaw(raw, r)
}
