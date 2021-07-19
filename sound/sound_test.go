package sound

import (
	"testing"
)

var path = "/Users/wellingtongadelha/Desktop/Musicas"

//go test -v -run ^TestNewSongs
func TestNewSongs(t *testing.T) {
	s := NewSongs()
	if s == nil {
		t.Errorf("TestNewSongs() return nil songs")
	}
}

//go test -v -run ^TestInitSongs
func TestInitSongs(t *testing.T) {
	s := NewSongs()
	err := s.Init(path)
	if err != nil {
		t.Errorf("TestInitSongs() return error: %v", err.Error())
	}
}

// go test -v -run ^TesGetDirs
func TestGetDirs(t *testing.T) {
	s := NewSongs()
	err := s.getDirs(path)
	if err != nil {
		t.Errorf("TestGetSongs() return error: %v", err.Error())
	}
}

// go test -v -run ^TestGetLibrary
func TestGetLibrary(t *testing.T) {
	s := NewSongs()
	s.Init(path)
	l := s.GetLibrary()
	if len(l) < 1 {
		t.Errorf("TestGetLibrary() return error: not values of list library")
	}
}

// go test -v -run ^TestGetSongs
func TestGetSongs(t *testing.T) {
	s := NewSongs()
	s.Init(path)
	l := s.GetSongs(1)
	if l == nil {
		t.Errorf("TestGetSongs() return error: list of value nil")
	}
}

// go test -v -run ^TestGetNamesLibrary
func TestGetNamesLibrary(t *testing.T) {
	s := NewSongs()
	s.Init(path)
	l := s.GetNames(s.GetLibrary())
	if l == nil {
		t.Errorf("TestGetNames() return error: list of library value nil")
	}
}

// go test -v -run ^TestGetNamesSongs
func TestGetNamesSongs(t *testing.T) {
	s := NewSongs()
	s.Init(path)
	l := s.GetNames(s.GetSongs(1))
	if l == nil {
		t.Errorf("TestGetNames() return error: list of songs value nil")
	}
}

// go test -v -run ^TestMusicPath
func TestMusicPath(t *testing.T) {
	s := NewSongs()
	s.Init(path)
	p := s.MusicPath(1, s.GetSongs(1))
	if p == "" {
		t.Errorf("TestMusicPath() return value empty.")
	}
}
