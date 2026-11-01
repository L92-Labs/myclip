package native

import "unicode/utf16"

// UTF-16LE, the encoding the Windows clipboard keeps text in. Pure and
// tested everywhere.

// utf16BytesToString decodes UTF-16LE clipboard bytes (dropping the NUL).
func utf16BytesToString(b []byte) string {
	if n := len(b); n >= 2 && b[n-1] == 0 && b[n-2] == 0 {
		b = b[:n-2]
	}
	u := make([]uint16, len(b)/2)
	for i := range u {
		u[i] = uint16(b[i*2]) | uint16(b[i*2+1])<<8
	}
	return string(utf16.Decode(u))
}

// utf16StringToBytes encodes a string as NUL-terminated UTF-16LE.
func utf16StringToBytes(s string) []byte {
	u := utf16.Encode([]rune(s))
	b := make([]byte, 2*len(u)+2)
	for i, v := range u {
		b[i*2], b[i*2+1] = byte(v), byte(v>>8)
	}
	return b
}
