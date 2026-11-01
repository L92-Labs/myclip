package native

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"math"
	"math/bits"
)

// Windows clipboard bitmaps arrive as CF_DIB: a BITMAPINFOHEADER (or the
// V4/V5 kind) followed by the palette and the pixels. These conversions
// are pure and tested everywhere.

const (
	biRGB       = 0
	biBitfields = 3
)

// dibToPNG decodes a CF_DIB payload — 24-bpp BI_RGB or 32-bpp BI_RGB or
// BI_BITFIELDS, top-down or bottom-up — into PNG bytes.
func dibToPNG(dib []byte) ([]byte, bool) {
	if len(dib) < 40 {
		return nil, false
	}
	size := int(binary.LittleEndian.Uint32(dib[0:4]))
	if size < 40 || size > len(dib) {
		return nil, false
	}
	w := int32(binary.LittleEndian.Uint32(dib[4:8]))
	h := int32(binary.LittleEndian.Uint32(dib[8:12]))
	bpp := binary.LittleEndian.Uint16(dib[14:16])
	compression := binary.LittleEndian.Uint32(dib[16:20])
	if w <= 0 || h == 0 || h == math.MinInt32 || compression != biRGB && compression != biBitfields {
		return nil, false // 0x80000000 stays negative when negated: reject it up front
	}
	topDown := h < 0
	if topDown {
		h = -h
	}
	if w > 1<<15 || h > 1<<15 {
		return nil, false // a corrupt header must not make a huge image
	}
	var rm, gm, bm, am maskInfo
	pixelsAt := size
	switch bpp {
	case 32:
		if compression == biBitfields {
			// The declared masks decide where each channel sits: standard
			// BGRA is only one layout of many. They sit 40 bytes into the
			// DIB either way — just after a 40-byte header, inside a
			// V4/V5 one.
			const maskAt = 40
			if size != 40 && size < maskAt+12 {
				return nil, false
			}
			if size == 40 {
				if len(dib) < 52 {
					return nil, false
				}
				pixelsAt = 52
			}
			at := maskAt
			r0 := binary.LittleEndian.Uint32(dib[at:])
			g0 := binary.LittleEndian.Uint32(dib[at+4:])
			b0 := binary.LittleEndian.Uint32(dib[at+8:])
			if r0|g0|b0 == 0 {
				return nil, false
			}
			var ok bool
			if rm, ok = maskLayout(r0); !ok {
				return nil, false
			}
			if gm, ok = maskLayout(g0); !ok {
				return nil, false
			}
			if bm, ok = maskLayout(b0); !ok {
				return nil, false
			}
			if size > 40 && at+16 <= len(dib) { // only a V4/V5 header declares alpha
				if a0 := binary.LittleEndian.Uint32(dib[at+12:]); a0 != 0 {
					am, ok = maskLayout(a0)
					if !ok {
						return nil, false
					}
				}
			}
		}
	case 24:
	default:
		return nil, false
	}
	stride := int((int(w)*int(bpp) + 31) / 32 * 4)
	if int64(pixelsAt)+int64(stride)*int64(h) > int64(len(dib)) {
		return nil, false
	}
	// BI_RGB's fourth byte is padding: unless some of it holds a real
	// alpha channel, the image is opaque. Declared masks say so directly.
	hasAlpha := am.bits > 0
	if bpp == 32 && compression == biRGB {
		for y := 0; y < int(h) && !hasAlpha; y++ {
			for x := 0; x < int(w); x++ {
				if dib[pixelsAt+y*stride+x*4+3] != 0 {
					hasAlpha = true
					break
				}
			}
		}
	}
	img := image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	for y := 0; y < int(h); y++ {
		row := dib[pixelsAt+y*stride:][:stride]
		src := y
		if !topDown { // bottom-up: the first row is the image's last
			src = int(h) - 1 - y
		}
		dst := img.Pix[src*img.Stride : src*img.Stride+int(w)*4]
		for x := 0; x < int(w); x++ {
			if bpp == 24 {
				dst[x*4], dst[x*4+1], dst[x*4+2], dst[x*4+3] = row[x*3+2], row[x*3+1], row[x*3], 255
				continue
			}
			var r, g, bl, a byte
			if compression == biBitfields {
				px := binary.LittleEndian.Uint32(row[x*4:])
				r, g, bl = maskByte(px, rm), maskByte(px, gm), maskByte(px, bm)
				a = 255
				if am.bits > 0 {
					a = maskByte(px, am)
				}
			} else {
				bl, g, r, a = row[x*4], row[x*4+1], row[x*4+2], row[x*4+3]
				if !hasAlpha {
					a = 255
				}
			}
			dst[x*4], dst[x*4+1], dst[x*4+2], dst[x*4+3] = r, g, bl, a
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, false
	}
	return out.Bytes(), true
}

// maskInfo is one channel's position in a 32-bit pixel.
type maskInfo struct{ shift, bits uint32 }

// maskLayout describes one channel's contiguous bit field, or rejects a
// mask that is not one (gaps, overlap checks happen across the masks).
func maskLayout(m uint32) (maskInfo, bool) {
	if m == 0 {
		return maskInfo{}, false
	}
	shift := uint32(bits.TrailingZeros32(m))
	width := uint32(32-bits.LeadingZeros32(m)) - shift
	if m>>shift != (1<<width)-1 {
		return maskInfo{}, false // not contiguous
	}
	return maskInfo{shift: shift, bits: width}, true
}

// maskByte scales one channel of a pixel to a full byte.
func maskByte(px uint32, m maskInfo) byte {
	v := (px >> m.shift) & ((1 << m.bits) - 1)
	if m.bits >= 8 {
		return byte(v >> (m.bits - 8))
	}
	return byte(v * 255 / ((1 << m.bits) - 1))
}

// pngToDIB encodes a PNG as a 32-bpp BI_RGB CF_DIB payload, bottom-up.
func pngToDIB(pngData []byte) ([]byte, bool) {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, false
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h > 1<<15 || w > 1<<15 {
		return nil, false
	}
	rgba, ok := img.(*image.NRGBA)
	if !ok {
		r := image.NewNRGBA(b)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				r.Set(x, y, img.At(x, y))
			}
		}
		rgba = r
	}
	stride := w * 4
	dib := make([]byte, 40+stride*h)
	binary.LittleEndian.PutUint32(dib[0:4], 40)         // biSize
	binary.LittleEndian.PutUint32(dib[4:8], uint32(w))  // biWidth
	binary.LittleEndian.PutUint32(dib[8:12], uint32(h)) // biHeight (bottom-up)
	binary.LittleEndian.PutUint16(dib[12:14], 1)        // biPlanes
	binary.LittleEndian.PutUint16(dib[14:16], 32)       // biBitCount
	binary.LittleEndian.PutUint32(dib[16:20], biRGB)    // biCompression
	binary.LittleEndian.PutUint32(dib[20:24], uint32(stride*h))
	for y := 0; y < h; y++ {
		src := rgba.Pix[(b.Min.Y+h-1-y)*rgba.Stride+b.Min.X*4:][:w*4]
		dst := dib[40+y*stride:]
		for x := 0; x < w; x++ {
			r, g, bl, a := src[x*4], src[x*4+1], src[x*4+2], src[x*4+3]
			if a == 0 {
				r, g, bl = 255, 255, 255 // fully transparent lands on white
			}
			dst[x*4], dst[x*4+1], dst[x*4+2], dst[x*4+3] = bl, g, r, a
		}
	}
	return dib, true
}
