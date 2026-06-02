package audiometa

// ParseTags detects FLAC, MP4/M4A (incl. AAC/ALAC in MP4), or MP3 and returns title, artist, album.
func ParseTags(buf []byte) (title, artist, album string) {
	if len(buf) >= 4 && string(buf[0:4]) == "fLaC" {
		t, a, b := ParseFLACTags(buf)
		if t != "" || a != "" || b != "" {
			return t, a, b
		}
		// Rare: ID3 appended after FLAC audio
		return ParseMP3Tags(buf)
	}
	if len(buf) >= 8 && string(buf[4:8]) == "ftyp" {
		t, a, b := ParseMP4Tags(buf)
		if t != "" || a != "" || b != "" {
			return t, a, b
		}
	}
	return ParseMP3Tags(buf)
}
