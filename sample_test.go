package impulse

import (
	"bytes"
	"testing"
)

var squareITS = []byte("IMPSsquare.wav\x00\x00\x00\x01\x01\x02square.wav\x00" +
	"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x83 " +
	"\x00\x00\x00\x04\x00\x00\x00\x05\x00\x00\x00\xab \x00\x00\x06\x00\x00" +
	"\x00\a\x00\x00\x00P\x00\x00\x00\b\t\n\x03\x80\x80\x80\x80\x80\x80\x80" +
	"\x80\x80\x80\x80\x80\x80\x80\x80\x80\u007f\u007f\u007f\u007f\u007f" +
	"\u007f\u007f\u007f\u007f\u007f\u007f\u007f\u007f\u007f\u007f\u007f")

func TestReadSample(t *testing.T) {
	// test invalid read on empty data
	r := bytes.NewReader([]byte{})
	if _, err := ReadSample(r); err == nil {
		t.Errorf("ReadSample() did not return error for empty data")
	}

	// test invalid read on bad data
	data := append([]byte("NOPE"), squareITS[4:]...)
	r = bytes.NewReader(data)
	if _, err := ReadSample(r); err == nil {
		t.Errorf("ReadSample() did not return error for bad data")
	}

	// test valid read
	r = bytes.NewReader(squareITS)
	s, err := ReadSample(r)
	if err != nil {
		t.Fatalf("ReadSample() returned error: %v", err)
	}

	// test fields
	if got, want := s.Filename, "square.wav"; got != want {
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
	if got, want := s.Name, "square.wav"; got != want {
		t.Errorf("Sample.Name == %#v; want %#v", got, want)
	}
	if got, want := s.Signed, true; got != want {
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
	want := []byte{128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
		128, 128, 128, 128, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127,
		127, 127, 127, 127, 127, 127}
	if got := s.Data; bytes.Compare(got, want) != 0 {
		t.Errorf("Sample.Data == %v; want %v", got, want)
	}

}
