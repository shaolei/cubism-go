package audioutils

import (
	"testing"
)

func TestDetectFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{"wav extension", "audio/sound.wav", "wav", false},
		{"wave extension", "audio/sound.wave", "wav", false},
		{"mp3 extension", "audio/sound.mp3", "mp3", false},
		{"wav uppercase", "audio/SOUND.WAV", "", true},
		{"ogg extension", "audio/sound.ogg", "", true},
		{"no extension", "audio/sound", "", true},
		{"path with dots", "audio.v2/sound.wav", "wav", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := DetectFormat(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFormat(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectFormat(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestDecodeAudioUnsupportedFormat(t *testing.T) {
	t.Parallel()

	_, _, err := DecodeAudio("flac", []byte{})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestNopCloser(t *testing.T) {
	t.Parallel()

	nc := NopCloser{}
	if err := nc.Close(); err != nil {
		t.Errorf("NopCloser.Close() returned error: %v", err)
	}
}
