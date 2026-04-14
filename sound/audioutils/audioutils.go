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
	default:
		err = fmt.Errorf("unsupported format: %s", format)
	}
	return
}
