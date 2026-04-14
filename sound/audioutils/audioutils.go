package audioutils

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/mewkiz/flac"
)

// nopCloser wraps an io.Reader to satisfy io.ReadCloser.
type NopCloser struct {
	io.Reader
}

func (NopCloser) Close() error { return nil }

// speakerInitMu ensures speaker.Init is called only once across all sound implementations.
var speakerInitMu sync.Mutex
var speakerInitialized bool

// InitSpeaker initializes the beep speaker with the given format.
// It is safe to call this from multiple goroutines; only the first call
// will actually initialize the speaker.
func InitSpeaker(format beep.Format) error {
	speakerInitMu.Lock()
	defer speakerInitMu.Unlock()
	if speakerInitialized {
		return nil
	}
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}
	speakerInitialized = true
	return nil
}

// DetectFormat returns the audio format name based on file extension.
func DetectFormat(fp string) (string, error) {
	ext := filepath.Ext(fp)
	switch ext {
	case ".wav", ".wave":
		return "wav", nil
	case ".mp3":
		return "mp3", nil
	case ".flac":
		return "flac", nil
	default:
		return "", fmt.Errorf("unsupported format: %s", ext)
	}
}

// DecodeAudio decodes audio data from the given buffer based on format.
func DecodeAudio(format string, buf []byte) (streamer beep.StreamSeekCloser, format2 beep.Format, err error) {
	switch format {
	case "wav":
		streamer, format2, err = wav.Decode(bytes.NewReader(buf))
	case "mp3":
		streamer, format2, err = mp3.Decode(NopCloser{bytes.NewReader(buf)})
	case "flac":
		streamer, format2, err = decodeFLAC(buf)
	default:
		err = fmt.Errorf("unsupported format: %s", format)
	}
	return
}

// flacStreamer implements beep.StreamSeekCloser for FLAC decoded PCM data
type flacStreamer struct {
	samples   [][2]float64 // interleaved stereo samples
	pos       int
	sampleRate beep.SampleRate
}

// decodeFLAC decodes FLAC audio data into a beep StreamSeekCloser
func decodeFLAC(buf []byte) (beep.StreamSeekCloser, beep.Format, error) {
	stream, err := flac.New(bytes.NewReader(buf))
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to parse FLAC: %w", err)
	}
	defer stream.Close()

	info := stream.Info
	sampleRate := beep.SampleRate(info.SampleRate)
	nChannels := int(info.NChannels)
	bitsPerSample := int(info.BitsPerSample)

	// Decode all frames and convert to float64 samples
	var allSamples [][2]float64
	maxVal := float64(int32(1) << (bitsPerSample - 1))

	for {
		frame, err := stream.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, beep.Format{}, fmt.Errorf("failed to parse FLAC frame: %w", err)
		}

		nSamples := len(frame.Subframes[0].Samples)
		for i := 0; i < nSamples; i++ {
			var s [2]float64
			// Channel 0
			s[0] = float64(frame.Subframes[0].Samples[i]) / maxVal
			// Channel 1 (or duplicate mono)
			if nChannels > 1 {
				s[1] = float64(frame.Subframes[1].Samples[i]) / maxVal
			} else {
				s[1] = s[0]
			}
			allSamples = append(allSamples, s)
		}
	}

	fs := &flacStreamer{
		samples:    allSamples,
		sampleRate: sampleRate,
	}

	format := beep.Format{
		SampleRate:  sampleRate,
		NumChannels: 2, // always output stereo for beep compatibility
		Precision:   bitsPerSample,
	}

	return fs, format, nil
}

func (fs *flacStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if fs.pos >= len(fs.samples) {
		return 0, false
	}
	n = copy(samples, fs.samples[fs.pos:])
	fs.pos += n
	return n, true
}

func (fs *flacStreamer) Err() error {
	return nil
}

func (fs *flacStreamer) Len() int {
	return len(fs.samples)
}

func (fs *flacStreamer) Position() int {
	return fs.pos
}

func (fs *flacStreamer) Seek(p int) error {
	if p < 0 || p > len(fs.samples) {
		return fmt.Errorf("flac: seek position %d out of range [0, %d]", p, len(fs.samples))
	}
	fs.pos = p
	return nil
}

func (fs *flacStreamer) Close() error {
	return nil
}
