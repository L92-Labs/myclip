//go:build darwin

package native

import (
	"errors"
	"strings"
	"sync"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// ErrNoOCR says this system cannot read text in images.
var ErrNoOCR = errors.New("native: text recognition is not available")

var (
	visionOnce sync.Once
	visionErr  error
)

func loadVision() error {
	visionOnce.Do(func() {
		if _, err := purego.Dlopen("/System/Library/Frameworks/Vision.framework/Vision", purego.RTLD_GLOBAL|purego.RTLD_NOW); err != nil {
			visionErr = errors.Join(ErrNoOCR, err)
			return
		}
		if objc.GetClass("VNRecognizeTextRequest") == 0 {
			visionErr = errors.Join(ErrNoOCR, errors.New("Vision text recognition needs macOS 10.15"))
		}
	})
	if visionErr == nil {
		visionErr = load()
	}
	return visionErr
}

// RecognizeText reads the text in an image (PNG, JPEG, TIFF, HEIC) with the
// Vision framework, on this Mac, line by line. It may take a moment: call it
// off the main thread.
func RecognizeText(image []byte) (string, error) {
	if err := loadVision(); err != nil {
		return "", err
	}
	var lines []string
	var err error
	withPool(func() {
		s := objc.RegisterName
		handler := objc.ID(objc.GetClass("VNImageRequestHandler")).Send(s("alloc")).
			Send(s("initWithData:options:"), nsData(image), objc.ID(objc.GetClass("NSDictionary")).Send(s("dictionary")))
		req := objc.ID(objc.GetClass("VNRecognizeTextRequest")).Send(s("alloc")).Send(s("init"))
		defer handler.Send(s("release"))
		defer req.Send(s("release"))
		req.Send(s("setRecognitionLevel:"), 0) // accurate
		req.Send(s("setUsesLanguageCorrection:"), true)
		requests := objc.ID(objc.GetClass("NSArray")).Send(s("arrayWithObject:"), req)
		if objc.Send[uint8](handler, s("performRequests:error:"), requests, uintptr(0)) == 0 {
			err = errors.New("native: Vision could not read the image")
			return
		}
		results := req.Send(s("results"))
		n := objc.Send[uint](results, s("count"))
		for i := uint(0); i < n; i++ {
			best := results.Send(selObjectAtIndex, i).Send(s("topCandidates:"), uint(1))
			if objc.Send[uint](best, s("count")) == 0 {
				continue
			}
			if line := strings.TrimSpace(goString(best.Send(selObjectAtIndex, uint(0)).Send(s("string")))); line != "" {
				lines = append(lines, line)
			}
		}
	})
	return strings.Join(lines, "\n"), err
}
