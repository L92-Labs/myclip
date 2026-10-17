package main

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/tiff"

	"github.com/L92-Labs/myclip/internal/darwin"
	"github.com/L92-Labs/myclip/internal/history"
)

// imageCache decodes history images on goroutines and keeps them, with
// apps' icons and names. Its methods run on the main thread; decoded
// images come back through invalidate.
type imageCache struct {
	invalidate func()

	mu      sync.Mutex
	thumbs  map[thumbKey]*ui.Bitmap
	full    map[*history.Item]*ui.Bitmap
	pending map[any]bool
	icons   map[string]*ui.Bitmap
	names   map[string]string
}

type thumbKey struct {
	it     *history.Item
	height int
}

func newImageCache(invalidate func()) *imageCache {
	return &imageCache{
		invalidate: invalidate,
		thumbs:     map[thumbKey]*ui.Bitmap{},
		full:       map[*history.Item]*ui.Bitmap{},
		pending:    map[any]bool{},
		icons:      map[string]*ui.Bitmap{},
		names:      map[string]string{},
	}
}

// screenScale is the sharpest display's scale, for pixels per DIP.
func screenScale() float64 {
	scale := 1.0
	for _, d := range mygo.Screen.Displays() {
		scale = max(scale, d.ScaleFactor)
	}
	return scale
}

// thumbnail is the item's image at most maxHeight DIPs high, or nil while
// it decodes or when it has none Go can read (HEIC).
func (c *imageCache) thumbnail(it *history.Item, maxHeight int) *ui.Bitmap {
	k := thumbKey{it, maxHeight}
	c.mu.Lock()
	defer c.mu.Unlock()
	if b, ok := c.thumbs[k]; ok || c.pending[k] {
		return b
	}
	data, _, ok := it.Image()
	if !ok {
		return nil
	}
	c.pending[k] = true
	px := int(float64(maxHeight) * screenScale())
	go c.decode(k, data, px, 340*screenScale(), func(b *ui.Bitmap) { c.thumbs[k] = b })
	return nil
}

// preview is the item's image for the preview pane, at most 1600 pixels.
func (c *imageCache) preview(it *history.Item) *ui.Bitmap {
	c.mu.Lock()
	defer c.mu.Unlock()
	if b, ok := c.full[it]; ok || c.pending[it] {
		return b
	}
	data, _, ok := it.Image()
	if !ok {
		return nil
	}
	c.pending[it] = true
	go c.decode(it, data, 1600, 1600, func(b *ui.Bitmap) { c.full[it] = b })
	return nil
}

// decode decodes data, scales it to fit maxH×maxW pixels, and stores the
// bitmap (nil when it does not decode) with put, under the lock.
func (c *imageCache) decode(k any, data []byte, maxH int, maxW float64, put func(*ui.Bitmap)) {
	var bmp *ui.Bitmap
	if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
		bmp = ui.NewBitmap(scaleDown(img, maxH, int(maxW)))
	}
	c.mu.Lock()
	delete(c.pending, k)
	put(bmp)
	c.mu.Unlock()
	mygo.RunOnMain(c.invalidate)
}

// scaleDown shrinks img to fit w×h pixels, keeping its proportions.
func scaleDown(img image.Image, h, w int) image.Image {
	b := img.Bounds()
	if b.Dy() <= h && b.Dx() <= w {
		return img
	}
	f := min(float64(h)/float64(b.Dy()), float64(w)/float64(b.Dx()))
	dst := image.NewRGBA(image.Rect(0, 0, max(1, int(float64(b.Dx())*f)), max(1, int(float64(b.Dy())*f))))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Src, nil)
	return dst
}

// forget drops the images of items no longer in the history.
func (c *imageCache) forget(items []*history.Item) {
	keep := make(map[*history.Item]bool, len(items))
	for _, it := range items {
		keep[it] = true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.thumbs {
		if !keep[k.it] {
			delete(c.thumbs, k)
		}
	}
	for it := range c.full {
		if !keep[it] {
			delete(c.full, it)
		}
	}
}

// appIcon is the icon of the app with bundleID, 16 DIPs, or nil.
func (c *imageCache) appIcon(bundleID string) *ui.Bitmap {
	if bundleID == "" {
		return nil
	}
	if b, ok := c.icons[bundleID]; ok {
		return b
	}
	var bmp *ui.Bitmap
	if path := darwin.AppPath(bundleID); path != "" {
		if data, err := mygo.Shell.FileIcon(path, int(16*screenScale())); err == nil {
			bmp, _ = ui.DecodeBitmap(data)
		}
	}
	c.icons[bundleID] = bmp
	return bmp
}

// appName is the app's name, as Finder shows it.
func (c *imageCache) appName(bundleID string) string {
	if n, ok := c.names[bundleID]; ok {
		return n
	}
	n := darwin.AppName(bundleID)
	c.names[bundleID] = n
	return n
}
