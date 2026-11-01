package native

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

// CF_HDROP is the clipboard/drag format for file paths: a DROPFILES
// header followed by a double-NUL-terminated path list, wide or not.
// Parsing and building it is pure and tested everywhere.

// dropFiles parses a CF_HDROP payload into file paths.
func dropFiles(b []byte) []string {
	if len(b) < 20 {
		return nil
	}
	offset := int(binary.LittleEndian.Uint32(b[0:4]))
	wide := binary.LittleEndian.Uint32(b[16:20]) != 0
	if offset < 20 || offset >= len(b) {
		return nil
	}
	list := b[offset:]
	var out []string
	if wide {
		for len(list) >= 2 {
			end := 0
			for end < len(list)/2 && le16(list, end) != 0 {
				end++
			}
			if end == 0 {
				break // the list's terminating NUL: done
			}
			out = append(out, string(utf16.Decode(le16s(list[:end*2]))))
			list = list[(end+1)*2:]
		}
		return out
	}
	for {
		end := bytes.IndexByte(list, 0)
		if end <= 0 {
			break
		}
		out = append(out, string(list[:end]))
		list = list[end+1:]
	}
	return out
}

// dropFilesBuild makes a wide-character CF_HDROP payload of paths.
func dropFilesBuild(files []string) []byte {
	var list []byte
	for _, f := range files {
		for _, v := range utf16.Encode([]rune(f)) {
			list = append(list, byte(v), byte(v>>8))
		}
		list = append(list, 0, 0) // each path ends in a NUL
	}
	list = append(list, 0, 0) // the list's double NUL
	hdr := make([]byte, 20)
	binary.LittleEndian.PutUint32(hdr[0:4], 20) // pFiles
	binary.LittleEndian.PutUint32(hdr[16:20], 1)
	return append(hdr, list...)
}

func le16(b []byte, i int) uint16 {
	return uint16(b[i*2]) | uint16(b[i*2+1])<<8
}

func le16s(b []byte) []uint16 {
	u := make([]uint16, len(b)/2)
	for i := range u {
		u[i] = le16(b, i)
	}
	return u
}
