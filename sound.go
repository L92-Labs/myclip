package main

import (
	"math"
	"sync"

	"github.com/egoist/mygo/audio"
)

// copySound is a short, soft pop drawn at launch: a sine falling from 1320
// to 880 Hz under a fast decay, 90 ms. No sound file ships with MyClip.
var copySound = sync.OnceValue(func() *audio.Sound {
	const rate, dur = 44100, 0.09
	n := int(rate * dur)
	samples := make([]float32, n)
	phase := 0.0
	for i := range samples {
		t := float64(i) / rate
		freq := 880 + 440*math.Exp(-t*40)
		phase += 2 * math.Pi * freq / rate
		attack := math.Min(1, t*400) // 2.5 ms, no click
		samples[i] = float32(0.25 * attack * math.Exp(-t*45) * math.Sin(phase))
	}
	return audio.NewSound(rate, 1, samples)
})

// playCopySound plays the pop when Play a sound on each copy is on.
func (a *app) playCopySound() {
	if a.settings.PlaySounds {
		copySound().Play()
	}
}
