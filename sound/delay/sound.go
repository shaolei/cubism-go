package delay

import (
	"os"

	"github.com/shaolei/cubism-go/sound"
	"github.com/shaolei/cubism-go/sound/audioutils"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Sound struct {
	fp       string
	streamer beep.StreamSeekCloser
	format   beep.Format
	ctrl     *beep.Ctrl
}

func LoadSound(fp string) (s sound.Sound, err error) {
	ds := &Sound{
		fp: fp,
	}
	return ds, nil
}

func (s *Sound) Decode() (err error) {
	if s.ctrl != nil {
		return
	}
	buf, err := os.ReadFile(s.fp)
	if err != nil {
		return
	}
	f, err := audioutils.DetectFormat(s.fp)
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
	if s.ctrl == nil {
		if err = s.Decode(); err != nil {
			return
		}
	}
	s.streamer.Seek(0)
	s.ctrl.Paused = false
	speaker.Play(s.ctrl)
	return
}

func (s *Sound) Close() {
	s.ctrl.Paused = true
	s.streamer.Seek(0)
}
