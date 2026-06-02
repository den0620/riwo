package audiometa

import (
	"encoding/binary"
	"testing"
)

func TestParseFLACVorbis(t *testing.T) {
	// fLaC + STREAMINFO (34) + VORBIS_COMMENT with TITLE=Hi
	var b []byte
	b = append(b, "fLaC"...)
	// Block 0: STREAMINFO, not last
	b = append(b, 0x00) // last=0, type=0
	b = append(b, 0x00, 0x00, 0x22)
	b = append(b, make([]byte, 34)...)
	// Block 1: VORBIS_COMMENT, last
	vb := buildVorbisComment("TITLE=Hi", "ARTIST=Me")
	b = append(b, 0x84) // last=1, type=4
	b = append(b, uint24be(len(vb))...)
	b = append(b, vb...)

	title, artist, album := ParseFLACTags(b)
	if title != "Hi" || artist != "Me" || album != "" {
		t.Fatalf("FLAC: got %q %q %q", title, artist, album)
	}
	title, artist, album = ParseTags(b)
	if title != "Hi" || artist != "Me" {
		t.Fatalf("ParseTags FLAC: got %q %q %q", title, artist, album)
	}
}

func uint24be(n int) []byte {
	return []byte{byte(n >> 16), byte(n >> 8), byte(n)}
}

func buildVorbisComment(pairs ...string) []byte {
	var out []byte
	// vendor string length + vendor (empty)
	out = binary.LittleEndian.AppendUint32(out, 0)
	// user comment list
	out = binary.LittleEndian.AppendUint32(out, uint32(len(pairs)))
	for _, p := range pairs {
		out = binary.LittleEndian.AppendUint32(out, uint32(len(p)))
		out = append(out, p...)
	}
	return out
}

func TestParseMP4Ilst(t *testing.T) {
	dataInner := make([]byte, 8+len("MyTitle"))
	copy(dataInner[8:], "MyTitle")
	dataBox := buildMP4Box("data", dataInner)
	namBox := buildMP4Box(tagName, dataBox)
	ilstBox := buildMP4Box("ilst", namBox)

	title, artist, album := ParseMP4Tags(ilstBox)
	if title != "MyTitle" || artist != "" || album != "" {
		t.Fatalf("ilst-only: %q %q %q", title, artist, album)
	}

	// Like a real file: ftyp then ilst (minimal, not strictly spec-complete)
	ftyp := make([]byte, 20)
	binary.BigEndian.PutUint32(ftyp[0:4], 20)
	copy(ftyp[4:8], "ftyp")
	copy(ftyp[8:12], "M4A ")
	binary.BigEndian.PutUint32(ftyp[12:16], 0)
	copy(ftyp[16:20], "M4A ")

	fileBuf := append(append([]byte{}, ftyp...), ilstBox...)
	title, artist, album = ParseTags(fileBuf)
	if title != "MyTitle" {
		t.Fatalf("ParseTags M4A: %q", title)
	}
}

func buildMP4Box(typ string, payload []byte) []byte {
	sz := 8 + len(payload)
	out := make([]byte, sz)
	binary.BigEndian.PutUint32(out[0:4], uint32(sz))
	copy(out[4:8], typ)
	copy(out[8:], payload)
	return out
}
