package sound

import (
	"bytes"
	"math"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ebitengine/oto/v3"
)

const (
	sampleRate     = 44100
	channelCount   = 2
	bytesPerSample = 2 // int16 LE

	clickVariantCount = 8
)

var (
	otoCtx   *oto.Context
	initOnce sync.Once

	clickVariants [clickVariantCount][]byte
	clickIdx      atomic.Uint32

	bufError     []byte
	bufNewline   []byte
	bufComplete  []byte
	bufPBest     []byte
	bufNavRow    []byte
	bufNavSelect []byte
)

// Init initialises the audio context and pre-generates all PCM buffers.
// Safe to call multiple times (sync.Once). Call it in a goroutine at startup.
func Init() {
	initOnce.Do(func() {
		ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
			SampleRate:   sampleRate,
			ChannelCount: channelCount,
			Format:       oto.FormatSignedInt16LE,
		})
		if err != nil {
			return
		}
		<-ready
		otoCtx = ctx

		// Pre-generate 8 mechanical click variants with slight pitch/amp variation.
		// Picking randomly prevents the robotic "identical click" effect at speed.
		for i := range clickVariantCount {
			bodyFreq := 380 + float64(i)*18 + rand.Float64()*20 // 380–560 Hz range
			amp := 0.22 + rand.Float64()*0.06
			clickVariants[i] = mechanicalClick(bodyFreq, 40*time.Millisecond, amp)
		}

		// Error: heavy thud — more body, less crack, lower frequency
		bufError = mechanicalThud(160, 70*time.Millisecond, 0.45)

		// Newline: G4→E4 two-tone descent
		bufNewline = append(
			sine(392, 14*time.Millisecond, 0.10),
			sine(330, 14*time.Millisecond, 0.10)...,
		)
		// Nav: G3 (subtle row hop) / C5 (confirmatory select click)
		bufNavRow = sine(196, 14*time.Millisecond, 0.09)
		bufNavSelect = sine(523, 16*time.Millisecond, 0.11)
		// Complete: C5→E5→G5→C6 major triad rise
		bufComplete = arpeggio([]note{
			{523, 65 * time.Millisecond},
			{659, 65 * time.Millisecond},
			{784, 65 * time.Millisecond},
			{1047, 160 * time.Millisecond},
		}, 0.28)
		// PersonalBest: same triad + E6 peak, louder and longer
		bufPBest = arpeggio([]note{
			{523, 50 * time.Millisecond},
			{659, 50 * time.Millisecond},
			{784, 50 * time.Millisecond},
			{1047, 60 * time.Millisecond},
			{1319, 210 * time.Millisecond},
		}, 0.38)
	})
}

func PlayCorrect() {
	// Round-robin through variants so no two consecutive clicks sound identical.
	i := clickIdx.Add(1) % clickVariantCount
	play(clickVariants[i])
}
func PlayError()        { play(bufError) }
func PlayNewline()      { play(bufNewline) }
func PlayComplete()     { play(bufComplete) }
func PlayPersonalBest() { play(bufPBest) }
func PlayNavRow()       { play(bufNavRow) }
func PlayNavSelect()    { play(bufNavSelect) }

func play(buf []byte) {
	if otoCtx == nil || len(buf) == 0 {
		return
	}
	go func() {
		p := otoCtx.NewPlayer(bytes.NewReader(buf))
		p.Play()
		for p.IsPlaying() {
			time.Sleep(time.Millisecond)
		}
		_ = p.Close()
	}()
}

// mechanicalClick synthesises a mechanical keyboard keypress:
//   - White-noise burst → the sharp "crack" of the switch actuating
//   - Body resonance sine → the "thock" from the keycap/plate vibrating
//
// Both use exponential decay (instant attack, ~8ms half-life) which is the
// signature of a physical impact — very different from a sine-wave tone.
func mechanicalClick(bodyFreq float64, dur time.Duration, amp float64) []byte {
	n := int(float64(sampleRate) * dur.Seconds())
	buf := make([]byte, n*channelCount*bytesPerSample)
	// 8ms half-life: loud transient that dies fast, like hitting a key
	decayConst := math.Log(2) / 0.008

	for i := range n {
		t := float64(i) / float64(sampleRate)
		env := math.Exp(-decayConst * t)

		noise := (rand.Float64()*2 - 1) * 0.70 // broadband crack
		body := math.Sin(2*math.Pi*bodyFreq*t) * 0.30 // resonant thock

		v := (noise + body) * amp * env
		// soft clip to prevent harsh clipping artifacts
		v = math.Tanh(v)
		s := int16(v * 32767)
		off := i * channelCount * bytesPerSample
		buf[off+0] = byte(s)
		buf[off+1] = byte(s >> 8)
		buf[off+2] = byte(s)
		buf[off+3] = byte(s >> 8)
	}
	return buf
}

// mechanicalThud is a heavier variant for errors — slower decay, lower body
// frequency, noise ratio inverted (more body, less crack).
func mechanicalThud(bodyFreq float64, dur time.Duration, amp float64) []byte {
	n := int(float64(sampleRate) * dur.Seconds())
	buf := make([]byte, n*channelCount*bytesPerSample)
	decayConst := math.Log(2) / 0.018 // 18ms half-life — slower, heavier

	for i := range n {
		t := float64(i) / float64(sampleRate)
		env := math.Exp(-decayConst * t)

		noise := (rand.Float64()*2 - 1) * 0.30
		body := math.Sin(2*math.Pi*bodyFreq*t) * 0.70

		v := (noise + body) * amp * env
		v = math.Tanh(v)
		s := int16(v * 32767)
		off := i * channelCount * bytesPerSample
		buf[off+0] = byte(s)
		buf[off+1] = byte(s >> 8)
		buf[off+2] = byte(s)
		buf[off+3] = byte(s >> 8)
	}
	return buf
}

// sine generates a sine wave with attack/release envelope.
func sine(freq float64, dur time.Duration, amp float64) []byte {
	n := int(float64(sampleRate) * dur.Seconds())
	buf := make([]byte, n*channelCount*bytesPerSample)
	attack := n / 5
	release := n / 3
	for i := range n {
		t := float64(i) / float64(sampleRate)
		env := 1.0
		switch {
		case i < attack:
			env = float64(i) / float64(attack)
		case i > n-release:
			env = float64(n-i) / float64(release)
		}
		s := int16(math.Sin(2*math.Pi*freq*t) * amp * env * 32767)
		off := i * channelCount * bytesPerSample
		buf[off+0] = byte(s)
		buf[off+1] = byte(s >> 8)
		buf[off+2] = byte(s)
		buf[off+3] = byte(s >> 8)
	}
	return buf
}

type note struct {
	freq float64
	dur  time.Duration
}

// arpeggio concatenates a sequence of sine-wave notes.
func arpeggio(notes []note, amp float64) []byte {
	var all []byte
	for _, n := range notes {
		all = append(all, sine(n.freq, n.dur, amp)...)
	}
	return all
}
