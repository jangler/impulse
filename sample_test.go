package impulse

import (
	"bytes"
	"testing"
)

func TestReadSample(t *testing.T) {
	// test valid read
	r := bytes.NewReader(MustAsset("_data/sine.its"))
	s, err := ReadSample(r)
	if err != nil {
		t.Errorf("readSample() returned error: %v", err)
	}

	// test fields
	if got, want := s.Filename, "sine.wav"; got != want {
		t.Errorf("Sample.Filename == %#v; want %#v", got, want)
	}
	if got, want := s.GlobalVolume, uint8(1); got != want {
		t.Errorf("Sample.GlobalVolume == %v; want %v", got, want)
	}
	if got, want := s.Flags, SampleAssociatedWithHeader; got != want {
		t.Errorf("Sample.Flags == %v; want %v", got, want)
	}
	if got, want := s.DefaultVolume, uint8(2); got != want {
		t.Errorf("Sample.DefaultVolume == %v; want %v", got, want)
	}
	if got, want := s.Name, "sine.wav"; got != want {
		t.Errorf("Sample.Name == %#v; want %#v", got, want)
	}
	if got, want := s.Signed, false; got != want {
		t.Errorf("Sample.Signed == %v; want %v", got, want)
	}
	if got, want := s.DefaultPan, uint8(3); got != want {
		t.Errorf("Sample.DefaultPan == %v; want %v", got, want)
	}
	if got, want := s.DefaultPanOn, true; got != want {
		t.Errorf("Sample.DefaultPanOn == %v; want %v", got, want)
	}
	if got, want := s.Length, uint32(32); got != want {
		t.Errorf("Sample.Length == %v; want %v", got, want)
	}
	if got, want := s.LoopBegin, uint32(4); got != want {
		t.Errorf("Sample.LoopBegin == %v; want %v", got, want)
	}
	if got, want := s.LoopEnd, uint32(5); got != want {
		t.Errorf("Sample.LoopEnd == %v; want %v", got, want)
	}
	if got, want := s.Speed, uint32(8363); got != want {
		t.Errorf("Sample.Speed == %v; want %v", got, want)
	}
	if got, want := s.SustainLoopBegin, uint32(6); got != want {
		t.Errorf("Sample.SustainLoopBegin == %v; want %v", got, want)
	}
	if got, want := s.SustainLoopEnd, uint32(7); got != want {
		t.Errorf("Sample.SustainLoopEnd == %v; want %v", got, want)
	}
	if got, want := s.VibratoSpeed, uint8(8); got != want {
		t.Errorf("Sample.VibratoSpeed == %v; want %v", got, want)
	}
	if got, want := s.VibratoDepth, uint8(9); got != want {
		t.Errorf("Sample.VibratoDepth == %v; want %v", got, want)
	}
	if got, want := s.VibratoRate, uint8(10); got != want {
		t.Errorf("Sample.VibratoRate == %v; want %v", got, want)
	}
	if got, want := s.VibratoWaveform, Random; got != want {
		t.Errorf("Sample.VibratoWaveform == %v; want %v", got, want)
	}

	// test invalid read on empty data
	r = bytes.NewReader([]byte{})
	if _, err := ReadSample(r); err == nil {
		t.Errorf("readSample() did not return error for empty data")
	}

	// test invalid read on bad data
	data := append([]byte("NOPE"), MustAsset("_data/sine.its")[4:]...)
	r = bytes.NewReader(data)
	if _, err := ReadSample(r); err == nil {
		t.Errorf("readSample() did not return error for bad data")
	}
}