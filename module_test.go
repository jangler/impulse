package impulse

import (
	"bytes"
	"testing"
)

var testIT = []byte("IMPMsong name\x00\x00\x00\x00\x00\x00\x00\x00" +
	"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x04\x10\x02\x00\x00\x00\x01\x00" +
	"\x01\x00\"\x18\x14\x02I\x00\x06\x00\x800\x06}\x80\f\x00\x00\x00\x00\x00" +
	"\x00\x00\x00\x00\x00@?>=<;:9876543210/.-,+*)('&%$#\"! \x1f\x1e\x1d\x1c" +
	"\x1b\x1a\x19\x18\x17\x16\x15\x14\x13\x12\x11\x10\x0f\x0e\r\f\v\n\t\b\a" +
	"\x06\x05\x04\x03\x02\x01\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f" +
	"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%" +
	"&'()*+,-./0123456789:;<=>?@\x01\xff\xd4\x00\x00\x00$\x01\x00\x00\x01" +
	"\x00\xfaFF8\xb9\xac\a\x00IMPSsquare.wav\x00\x00\x00\x01\x01\x02square.w" +
	"av\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01" +
	"\x83 \x00\x00\x00\x04\x00\x00\x00\x05\x00\x00\x00\xab \x00\x00\x06\x00" +
	"\x00\x00\a\x00\x00\x00r\x01\x00\x00\b\t\n\x03F\x00 \x00\x00\x00\x00\x00" +
	"\x81\x034\x01\x00\x00\x81!2\x00\x00\x010\x00\x00\x012\x00\x00\x014\x00" +
	"\x00\x810\x00\x00\x01\x00\x00\x81\x01\xff\x00\x00\x81!2\x00\x00\x810" +
	"\x00\x00\x01\x00\x00\x81\x01\xff\x00\x00\x81!4\x00\x00\x017\x00\x00\x81" +
	"0\x00\x00\x81\x01\xff\x00\x00\x80\x80\x80\x80\x80\x80\x80\x80\x80\x80" +
	"\x80\x80\x80\x80\x80\x80\u007f\u007f\u007f\u007f\u007f\u007f\u007f" +
	"\u007f\u007f\u007f\u007f\u007f\u007f\u007f\u007f\u007f")

func TestModule(t *testing.T) {
	// test invalid read on empty data
	r := bytes.NewReader([]byte{})
	if _, err := ReadModule(r); err == nil {
		t.Errorf("ReadModule() did not return error for empty data")
	}

	// test invalid read on bad data
	data := append([]byte("NOPE"), testIT[4:]...)
	r = bytes.NewReader(data)
	if _, err := ReadModule(r); err == nil {
		t.Errorf("ReadModule() did not return error for bad data")
	}

	// test valid read
	r = bytes.NewReader(testIT)
	m, err := ReadModule(r)
	if err != nil {
		t.Fatalf("ReadModule() returned error: %v", err)
	}

	// test fields
	if got, want := m.SongName, "song name"; got != want {
		t.Errorf("Module.SongName == %#v; want %#v", got, want)
	}
	if got, want := m.GlobalVolume, uint8(128); got != want {
		t.Errorf("Module.GlobalVolume == %v; want %v", got, want)
	}
	if got, want := m.MixingVolume, uint8(48); got != want {
		t.Errorf("Module.MixingVolume == %v; want %v", got, want)
	}
	if got, want := m.InitialSpeed, uint8(6); got != want {
		t.Errorf("Module.InitialSpeed == %v; want %v", got, want)
	}
	if got, want := m.InitialTempo, uint8(125); got != want {
		t.Errorf("Module.InitialTempo == %v; want %v", got, want)
	}
	if got, want := m.Separation, uint8(128); got != want {
		t.Errorf("Module.Separation == %v; want %v", got, want)
	}
	if got, want := m.PitchWheelDepth, uint8(12); got != want {
		t.Errorf("Module.PitchWheelDepth == %v; want %v", got, want)
	}
	for i, v := range m.ChannelPanning {
		if got, want := v, uint8(64-i); got != want {
			t.Errorf("Module.ChannelPanning[%d] == %v; want %v", i, got, want)
		}
	}
	for i, v := range m.ChannelVolume {
		if got, want := v, uint8(i+1); got != want {
			t.Errorf("Module.ChannelVolume[%d] == %v; want %v", i, got, want)
		}
	}
	want := []byte{1, 255}
	if got := m.OrderList; bytes.Compare(got, want) != 0 {
		t.Errorf("Module.OrderList == %v; want %v", got, want)
	}
}
