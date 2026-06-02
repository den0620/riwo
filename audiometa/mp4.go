package audiometa

import (
	"encoding/binary"
	"strings"
)

// iTunes-style QuickTime tag types (4-char codes).
var (
	tagName   = string([]byte{0xa9, 'n', 'a', 'm'})
	tagArtist = string([]byte{0xa9, 'A', 'R', 'T'})
	tagAlbum  = string([]byte{0xa9, 'a', 'l', 'b'})
	tagAlbumArtist = "aART"
)

// ParseMP4Tags walks MP4/M4A boxes and reads ©nam, ©ART, ©alb from an ilst atom.
func ParseMP4Tags(buf []byte) (title, artist, album string) {
	var t, a, b string
	var walk func([]byte)
	walk = func(data []byte) {
		off := 0
		for off+8 <= len(data) {
			boxEnd, typ, payload, ok := nextMP4Box(data, off)
			if !ok {
				break
			}
			switch typ {
			case "moov", "trak", "mdia", "minf", "stbl", "edts", "dinf", "clip", "tapt", "udta", "cmov", "mvex":
				walk(payload)
			case "meta":
				if len(payload) >= 4 {
					walk(payload[4:])
				}
			case "ilst":
				parseIlst(payload, func(atomType, val string) {
					switch atomType {
					case tagName:
						if t == "" {
							t = val
						}
					case tagArtist:
						if a == "" {
							a = val
						}
					case tagAlbum:
						if b == "" {
							b = val
						}
					case tagAlbumArtist:
						if a == "" {
							a = val
						}
					}
				})
			default:
				// ftyp, mdat, free, hdlr, ...
			}
			off = boxEnd
		}
	}
	walk(buf)
	return strings.TrimSpace(t), strings.TrimSpace(a), strings.TrimSpace(b)
}

func parseIlst(data []byte, set func(atomType, val string)) {
	off := 0
	for off+8 <= len(data) {
		boxEnd, typ, payload, ok := nextMP4Box(data, off)
		if !ok {
			break
		}
		if txt := extractDataAtomText(payload); txt != "" {
			set(typ, txt)
		}
		off = boxEnd
	}
}

// extractDataAtomText finds the first 'data' child and decodes its payload.
func extractDataAtomText(itemPayload []byte) string {
	off := 0
	for off+8 <= len(itemPayload) {
		boxEnd, typ, payload, ok := nextMP4Box(itemPayload, off)
		if !ok {
			break
		}
		if typ == "data" {
			return decodeMP4DataValue(payload)
		}
		// some atoms nest mean/name (----); not handled
		off = boxEnd
	}
	return ""
}

func decodeMP4DataValue(fullBoxPayload []byte) string {
	if len(fullBoxPayload) < 8 {
		return ""
	}
	// FullBox: version (1) + flags (3) + NULL (4), then UTF-8 text for type 1.
	body := fullBoxPayload[8:]
	return strings.TrimSpace(string(trimBOMUTF8(body)))
}

func trimBOMUTF8(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		return b[3:]
	}
	return b
}

// nextMP4Box returns the end offset of this box, the type, and the payload after the box header.
func nextMP4Box(data []byte, off int) (boxEnd int, typ string, payload []byte, ok bool) {
	if off+8 > len(data) {
		return 0, "", nil, false
	}
	sz := binary.BigEndian.Uint32(data[off : off+4])
	typ = string(data[off+4 : off+8])
	headerLen := 8
	boxLen := int(sz)
	if sz == 1 {
		if off+16 > len(data) {
			return 0, "", nil, false
		}
		boxLen = int(binary.BigEndian.Uint64(data[off+8 : off+16]))
		headerLen = 16
	} else if sz == 0 {
		boxLen = len(data) - off
	}
	if boxLen < headerLen || off+boxLen > len(data) {
		return 0, "", nil, false
	}
	payload = data[off+headerLen : off+boxLen]
	boxEnd = off + boxLen
	return boxEnd, typ, payload, true
}
