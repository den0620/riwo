package audiometa

import (
	"encoding/binary"
	"strings"
)

// ParseFLACTags reads Vorbis comment block (metadata block type 4) from a FLAC stream.
func ParseFLACTags(buf []byte) (title, artist, album string) {
	if len(buf) < 4 || string(buf[0:4]) != "fLaC" {
		return "", "", ""
	}
	pos := 4
	for pos+4 <= len(buf) {
		hdr := buf[pos]
		last := hdr&0x80 != 0
		btype := hdr & 0x7f
		length := int(buf[pos+1])<<16 | int(buf[pos+2])<<8 | int(buf[pos+3])
		pos += 4
		if length < 0 || pos+length > len(buf) {
			return "", "", ""
		}
		block := buf[pos : pos+length]
		if btype == 4 {
			return parseVorbisCommentBlock(block)
		}
		pos += length
		if last {
			break
		}
	}
	return "", "", ""
}

func parseVorbisCommentBlock(data []byte) (title, artist, album string) {
	if len(data) < 8 {
		return "", "", ""
	}
	pos := 0
	vendorLen := int(binary.LittleEndian.Uint32(data[0:4]))
	pos += 4
	if vendorLen < 0 || pos+vendorLen > len(data) {
		return "", "", ""
	}
	pos += vendorLen
	if pos+4 > len(data) {
		return "", "", ""
	}
	n := int(binary.LittleEndian.Uint32(data[pos : pos+4]))
	pos += 4
	for i := 0; i < n && pos+4 <= len(data); i++ {
		clen := int(binary.LittleEndian.Uint32(data[pos : pos+4]))
		pos += 4
		if clen < 0 || pos+clen > len(data) {
			break
		}
		line := string(data[pos : pos+clen])
		pos += clen
		key, val, ok := splitVorbisField(line)
		if !ok {
			continue
		}
		switch strings.ToUpper(key) {
		case "TITLE":
			if title == "" {
				title = val
			}
		case "ARTIST":
			if artist == "" {
				artist = val
			}
		case "ALBUM":
			if album == "" {
				album = val
			}
		}
	}
	return strings.TrimSpace(title), strings.TrimSpace(artist), strings.TrimSpace(album)
}

func splitVorbisField(line string) (key, val string, ok bool) {
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			return line[:i], strings.TrimSpace(line[i+1:]), true
		}
	}
	return "", "", false
}
