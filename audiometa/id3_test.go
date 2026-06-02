package audiometa

import "testing"

func TestParseID3v1Only(t *testing.T) {
	block := make([]byte, 128)
	copy(block[0:], "TAG")
	copy(block[3:33], pad32("Song"))
	copy(block[33:63], pad32("Artist"))
	copy(block[63:93], pad32("Album"))

	title, artist, album := ParseMP3Tags(block)
	if title != "Song" || artist != "Artist" || album != "Album" {
		t.Fatalf("got %q / %q / %q", title, artist, album)
	}
}

func pad32(s string) []byte {
	out := make([]byte, 30)
	copy(out, s)
	return out
}
