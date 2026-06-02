package audiometa

import (
	"bytes"
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

// ParseMP3Tags extracts title, artist, and album from MP3 file bytes using ID3v2.3/v2.4 and ID3v1.
// For other formats use [ParseTags].
func ParseMP3Tags(buf []byte) (title, artist, album string) {
	if len(buf) < 10 {
		return "", "", ""
	}
	if string(buf[0:3]) == "ID3" {
		ver := buf[3]
		tagSize := syncsafe32(buf[6:10])
		dataEnd := 10 + tagSize
		if dataEnd > len(buf) {
			dataEnd = len(buf)
		}
		switch ver {
		case 3:
			title, artist, album = id3v23ParseFrames(buf[10:dataEnd])
		case 4:
			title, artist, album = id3v24ParseFrames(buf[10:dataEnd])
		}
	}
	if len(buf) >= 128 {
		if t, a, b := parseID3v1(buf[len(buf)-128:]); t != "" || a != "" || b != "" {
			if title == "" {
				title = t
			}
			if artist == "" {
				artist = a
			}
			if album == "" {
				album = b
			}
		}
	}
	return strings.TrimSpace(title), strings.TrimSpace(artist), strings.TrimSpace(album)
}

func syncsafe32(b []byte) int {
	if len(b) < 4 {
		return 0
	}
	return int(b[0]&0x7f)<<21 | int(b[1]&0x7f)<<14 | int(b[2]&0x7f)<<7 | int(b[3]&0x7f)
}

func id3v23ParseFrames(tag []byte) (title, artist, album string) {
	return id3v23or24Frames(tag, false)
}

func id3v24ParseFrames(tag []byte) (title, artist, album string) {
	return id3v23or24Frames(tag, true)
}

func id3v23or24Frames(tag []byte, v24 bool) (title, artist, album string) {
	pos := 0
	for pos+10 <= len(tag) {
		id := string(tag[pos : pos+4])
		if id == "\x00\x00\x00\x00" || id[0] == 0 {
			break
		}
		var sz int
		if v24 {
			sz = syncsafe32(tag[pos+4 : pos+8])
		} else {
			sz = int(binary.BigEndian.Uint32(tag[pos+4 : pos+8]))
		}
		if sz < 0 || pos+10+sz > len(tag) {
			break
		}
		raw := tag[pos+10 : pos+10+sz]
		switch id {
		case "TIT2":
			if title == "" {
				title = decodeID3TextFrame(raw)
			}
		case "TPE1":
			if artist == "" {
				artist = decodeID3TextFrame(raw)
			}
		case "TALB":
			if album == "" {
				album = decodeID3TextFrame(raw)
			}
		}
		pos += 10 + sz
	}
	return title, artist, album
}

func decodeID3TextFrame(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	enc := raw[0]
	payload := raw[1:]
	switch enc {
	case 0:
		return strings.TrimRight(string(payload), "\x00")
	case 3:
		return strings.TrimRight(string(payload), "\x00")
	case 1:
		return decodeUTF16BOMPayload(payload)
	case 2:
		return decodeUTF16BE(payload)
	default:
		return strings.TrimRight(string(payload), "\x00")
	}
}

func decodeUTF16BOMPayload(p []byte) string {
	if len(p) < 2 {
		return ""
	}
	off := 0
	bigEndian := true
	if p[0] == 0xfe && p[1] == 0xff {
		off = 2
		bigEndian = true
	} else if p[0] == 0xff && p[1] == 0xfe {
		off = 2
		bigEndian = false
	}
	return utf16ToString(p[off:], bigEndian)
}

func decodeUTF16BE(p []byte) string {
	return utf16ToString(p, true)
}

func utf16ToString(p []byte, bigEndian bool) string {
	if len(p) < 2 {
		return ""
	}
	u := make([]uint16, 0, len(p)/2)
	for i := 0; i+1 < len(p); i += 2 {
		var r uint16
		if bigEndian {
			r = uint16(p[i])<<8 | uint16(p[i+1])
		} else {
			r = uint16(p[i]) | uint16(p[i+1])<<8
		}
		u = append(u, r)
	}
	return strings.TrimRight(string(utf16.Decode(u)), "\x00")
}

func parseID3v1(block []byte) (title, artist, album string) {
	if len(block) < 128 || string(block[0:3]) != "TAG" {
		return "", "", ""
	}
	title = trimID3v1Field(block[3:33])
	artist = trimID3v1Field(block[33:63])
	album = trimID3v1Field(block[63:93])
	return title, artist, album
}

func trimID3v1Field(b []byte) string {
	s := string(bytes.TrimRight(b, "\x00"))
	return strings.TrimSpace(s)
}
