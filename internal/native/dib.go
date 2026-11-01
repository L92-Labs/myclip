package native

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
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
	if w <= 0 || h == 0 || compression != biRGB && compression != biBitfields {
		return nil, false
	}
	topDown := h < 0
	if topDown {
		h = -h
	}
	if w > 1<<15 || h > 1<<15 {
		return nil, false // a corrupt header must not make a huge image
	}
	pixelsAt := size
	switch bpp {
	case 32:
		// BI_BITFIELDS may carry three DWORD masks after the header.
		if compression == biBitfields && size == 40 && len(dib) >= size+12 {
			pixelsAt = size + 12
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
	// alpha channel, the image is opaque.
	hasAlpha := false
	if bpp == 32 {
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
			if bpp == 32 {
				bl, g, r, a := row[x*4], row[x*4+1], row[x*4+2], row[x*4+3]
				if !hasAlpha {
					a = 255
				}
				dst[x*4], dst[x*4+1], dst[x*4+2], dst[x*4+3] = r, g, bl, a
			} else {
				dst[x*4], dst[x*4+1], dst[x*4+2], dst[x*4+3] = row[x*3+2], row[x*3+1], row[x*3], 255
			}
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, false
	}
	return out.Bytes(), true
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
