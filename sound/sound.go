package sound

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Songs struct {
	Library []string
}

var supportedFormats = []string{".mp3", ".wav", ".flac"}
var mainCtrl *beep.Ctrl
var s beep.StreamSeekCloser
var format beep.Format
var volume = &effects.Volume{
	Base: 2,
}

//NewSongs return of instance of struct Songs
func NewSongs() *Songs {
	return &Songs{}
}

//Init function initialise library the songs.
func (s *Songs) Init(path string) error {
	err := s.getDirs(path)
	if err != nil {
		return err
	}
	return nil
}

//getDirs return slice of all dirs.
func (s *Songs) getDirs(path string) error {
	var m []string
	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {

		if info.IsDir() {
			m = append(m, dir)
		}
		return err
	})
	s.Library = m
	return err
}

//GetLibrary return of list the library directorys.
func (s *Songs) GetLibrary() []string {
	var l []string
	for _, v := range s.Library {
		l = append(l, v)
	}
	return l
}

//GetSongs return of list the songs the index library
func (s *Songs) GetSongs(index int) []string {
	var sl []string
	err := filepath.Walk(s.Library[index], func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && strings.HasSuffix(path, ".mp3") {
			sl = append(sl, path)
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	return sl
}

// GetNames return of base name from filepaths.
func (s *Songs) GetNames(names []string) []string {
	var n []string
	for _, v := range names {
		n = append(n, filepath.Base(v))
	}
	return n
}

// MusicPath return path of music from row selected
func (s *Songs) MusicPath(index int, songs []string) string {
	return songs[index]
}

//PlaySong execute file .mp3
func PlaySong(input string) (int, error) {
	f, err := os.Open(input)
	if err != nil {
		return 0, err
	}

	switch fileExt := filepath.Ext(input); fileExt {
	case ".mp3":
		s, format, err = mp3.Decode(f)
	case ".wav":
		s, format, err = wav.Decode(f)
	case ".flac":
		s, format, err = flac.Decode(f)
	}
	if err != nil {
		return 0, err
	}
	volume.Streamer = s
	mainCtrl = &beep.Ctrl{Streamer: volume}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(mainCtrl)
	return int(float32(s.Len()) / float32(format.SampleRate)), nil
}

func PauseSong(state bool) {
	speaker.Lock()
	mainCtrl.Paused = state
	speaker.Unlock()
}

func Seek(pos int) error {
	speaker.Lock()
	err := s.Seek(pos * int(format.SampleRate))
	speaker.Unlock()
	return err
}

func SetVolume(percent int) {
	if percent == 0 {
		volume.Silent = true
	} else {
		volume.Silent = false
		volume.Volume = -float64(100-percent) / 100.0 * 5
	}
}
