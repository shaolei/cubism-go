package normal

import (
	"os"

	"github.com/shaolei/cubism-go/sound"
	"github.com/shaolei/cubism-go/sound/audioutils"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Sound struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
	ctrl     *beep.Ctrl
}

func LoadSound(fp string) (s sound.Sound, err error) {
	ds := &Sound{}
	buf, err := os.ReadFile(fp)
	if err != nil {
		return
	}
	return ds, ds.Decode(fp, buf)
}

func (s *Sound) Decode(fp string, buf []byte) (err error) {
	if s.ctrl != nil {
		return
	}
	f, err := audioutils.DetectFormat(fp)
	if err != nil {
		return
	}
	s.streamer, s.format, err = audioutils.DecodeAudio(f, buf)
	if err != nil {
		return
	}
	s.ctrl = &beep.Ctrl{Streamer: s.streamer}
	err = audioutils.InitSpeaker(s.format)
	return
}

func (s *Sound) Play() (err error) {
	s.streamer.Seek(0)
	s.ctrl.Paused = false
	speaker.Play(s.ctrl)
	return
}

func (s *Sound) Close() {
	s.ctrl.Paused = true
	s.streamer.Seek(0)
}
