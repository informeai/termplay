package ui

import (
	"context"
	"fmt"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	sound "github.com/informeai/termplay/sound"
)

type Ui struct {
	mainContainer *termui.Grid
	shortCuts     *widgets.Paragraph
	library       *widgets.List
	music         *widgets.List
	currentTime   *widgets.Gauge
	volume        *widgets.Gauge
}

//NewUi return instance of Ui struct.
func NewUi() *Ui {
	return &Ui{}
}

// Init function initialize ui
func (u *Ui) Run(path string) {
	// init songs
	s := sound.NewSongs()
	err := s.Init(path)
	if err != nil {
		panic(err)
	}
	// init termui
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()
	u.createWidgets(s)
	u.runEvents(s)
}

//create widget
func (u *Ui) createWidgets(s *sound.Songs) {
	// terminal dimensions
	termWidth, termHeight := termui.TerminalDimensions()
	//Create widgets
	u.mainContainer = termui.NewGrid()
	u.mainContainer.SetRect(0, 0, termWidth, termHeight)
	//short cuts ui
	u.shortCuts = widgets.NewParagraph()
	u.shortCuts.Title = "Keys"
	u.shortCuts.TitleStyle.Fg = termui.ColorYellow
	u.shortCuts.Block.BorderStyle.Fg = termui.ColorMagenta
	u.shortCuts.Text = "[ p ](fg:green,bg:magenta)[ Play/Stop ](fg:magenta,bg:green) " +
		"[ Left ](fg:green,bg:magenta)[ Library ](fg:magenta,bg:green) " +
		"[ Right ](fg:green,bg:magenta)[ Musics ](fg:magenta,bg:green) " +
		"[ + ](fg:green,bg:magenta)[ +Volume ](fg:magenta,bg:green) " +
		"[ - ](fg:green,bg:magenta)[ -Volume ](fg:magenta,bg:green) " +
		"[ q ](fg:green,bg:red)[ Exit ](fg:magenta,bg:green) "
	u.shortCuts.Border = true
	//library ui
	u.library = widgets.NewList()
	u.library.Border = true
	u.library.Title = "Library"
	//get library
	u.library.Rows = s.GetNames(s.GetLibrary())
	// verify row selected
	u.library.SelectedRowStyle = termui.NewStyle(termui.ColorGreen, termui.ColorBlack)
	u.library.TitleStyle.Fg = termui.ColorGreen
	u.library.Block.BorderStyle.Fg = termui.ColorMagenta
	// music ui
	u.music = widgets.NewList()
	u.music.Title = "Musics"
	u.music.Border = true
	u.music.TitleStyle.Fg = termui.ColorCyan
	u.music.Block.BorderStyle.Fg = termui.ColorMagenta
	u.music.Rows = s.GetNames(s.GetSongs(u.library.SelectedRow))
	// time
	u.currentTime = widgets.NewGauge()
	u.currentTime.Title = "Stop"
	u.currentTime.TitleStyle.Fg = termui.ColorCyan
	u.currentTime.Label = "00:00/00:00"
	u.currentTime.LabelStyle.Fg = termui.ColorGreen
	u.currentTime.Percent = 0
	u.currentTime.BarColor = termui.ColorBlue
	u.currentTime.PaddingTop = 1
	u.currentTime.PaddingLeft = 1
	u.currentTime.PaddingRight = 2
	u.currentTime.Block.BorderStyle.Fg = termui.ColorMagenta
	// volume
	u.volume = widgets.NewGauge()
	u.volume.Title = "volume"
	u.volume.Percent = 80
	u.volume.Label = fmt.Sprint(u.volume.Percent, "%")
	u.volume.TitleStyle = termui.NewStyle(termui.ColorCyan)
	u.volume.Block.BorderStyle = termui.NewStyle(termui.ColorMagenta)
	u.volume.BarColor = termui.ColorBlue
	u.volume.LabelStyle = termui.NewStyle(termui.ColorGreen)
	u.volume.LabelStyle.Modifier = termui.Modifier(termui.AlignCenter)
	u.volume.PaddingTop = 1
	u.volume.PaddingLeft = 1
	u.volume.PaddingRight = 1
	sound.SetVolume(u.volume.Percent)

	// set mainContainer
	u.mainContainer.Set(
		termui.NewRow(1.0/9, u.shortCuts),
		termui.NewRow(
			1.0/1.3,
			termui.NewCol(1.0/3, u.library),
			termui.NewCol(1.0/1.5, u.music),
		),
		termui.NewRow(
			1.0/8,
			termui.NewCol(1.0/1.5, u.currentTime),
			termui.NewCol(1.0/3, u.volume),
		),
	)

	termui.Render(u.mainContainer)
}

func (u *Ui) runEvents(s *sound.Songs) {
	var stateLibrayMusic string = "library"
	var statePlay bool = false
	for e := range termui.PollEvents() {
		if e.Type == termui.KeyboardEvent {
			switch e.ID {
			case "q":
				return
			case "j", "<Down>":
				if stateLibrayMusic == "library" {
					u.library.ScrollDown()
					u.music.Rows = s.GetNames(s.GetSongs(u.library.SelectedRow))
				} else if statePlay == false {
					u.music.ScrollDown()
				}
			case "k", "<Up>":
				if stateLibrayMusic == "library" {
					u.library.ScrollUp()
					u.music.Rows = s.GetNames(s.GetSongs(u.library.SelectedRow))
				} else if statePlay == false {
					u.music.ScrollUp()
				}
			case "=", "+":
				if u.volume.Percent >= 0 && u.volume.Percent < 100 {
					u.volume.Percent += 10
					u.volume.Label = fmt.Sprint(u.volume.Percent, "%")
					sound.SetVolume(u.volume.Percent)

				}
			case "-", "_":
				if u.volume.Percent > 0 && u.volume.Percent <= 100 {
					u.volume.Percent -= 10
					u.volume.Label = fmt.Sprint(u.volume.Percent, "%")
					sound.SetVolume(u.volume.Percent)

				}
			case "<Left>":
				if statePlay == false {
					stateLibrayMusic = "library"
					u.library.ScrollTop()
					u.library.SelectedRowStyle = termui.NewStyle(termui.ColorGreen, termui.ColorBlack)
					u.music.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorBlack)
					u.music.Rows = s.GetNames(s.GetSongs(u.library.SelectedRow))
					u.music.ScrollTop()
				}
			case "<Right>":
				if statePlay == false {
					stateLibrayMusic = "music"
					u.library.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorBlack)
					u.music.SelectedRowStyle = termui.NewStyle(termui.ColorMagenta, termui.ColorBlack)
					u.music.Rows = s.GetNames(s.GetSongs(u.library.SelectedRow))
					u.music.ScrollTop()
				}
			case "p":
				actualIndex := u.music.SelectedRow
				lenghtMusics := len(u.music.Rows)
				pos := 0
				ticker := time.NewTicker(time.Second)
				_, cancel := context.WithCancel(context.Background())
				u.currentTime.Label = fmt.Sprintf("00:00 / 00:00")
				if stateLibrayMusic == "music" && statePlay == false {

					go func() {
						var nextIndex int = 0
						var songLen int = 0
						u.currentTime.Title = "Playing"

						path := s.MusicPath(actualIndex+nextIndex, s.GetSongs(u.library.SelectedRow))
						songLen, err := sound.PlaySong(path)
						if err != nil {
							panic(err)
						}

						for {
							time.Sleep(time.Second)
							pos++
							u.currentTime.Label = fmt.Sprintf("%d:%.2d / %d:%.2d", pos/60, pos%60, songLen/60, songLen%60)
							u.currentTime.Percent = int((pos * 100) / songLen)
							if statePlay == false {
								defer cancel()
								ticker.Stop()
								pos = 0
								songLen = 0
								u.currentTime.Title = "Stop"
								u.currentTime.Percent = 0
								u.currentTime.Label = fmt.Sprintf("00:00 / 00:00")
								u.music.SelectedRow = 0
								break

							} else if pos == songLen && statePlay == true && u.music.SelectedRow < lenghtMusics-1 {
								ticker.Stop()
								ticker.Reset(time.Second)
								pos = 0
								u.currentTime.Percent = 0
								u.currentTime.Label = fmt.Sprintf("00:00 / 00:00")
								nextIndex++
								u.music.SelectedRow = actualIndex + nextIndex
								//next song
								path := s.MusicPath(actualIndex+nextIndex, s.GetSongs(u.library.SelectedRow))
								songLen, err := sound.PlaySong(path)
								if err != nil {
									panic(err)
								}
								u.currentTime.Title = "Playing"
								u.currentTime.Label = fmt.Sprintf("%d:%.2d / %d:%.2d", pos/60, pos%60, songLen/60, songLen%60)
								u.currentTime.Percent = int((pos * 100) / songLen)
								termui.Render(u.mainContainer)
							} else if pos == songLen && statePlay == true && u.music.SelectedRow == lenghtMusics-1 {
								statePlay = false
								u.currentTime.Title = "Stop"
								u.currentTime.Percent = 0
								u.currentTime.Label = fmt.Sprintf("00:00 / 00:00")
								u.music.SelectedRow = 0
								termui.Render(u.mainContainer)
							}
							termui.Render(u.mainContainer)
						}
					}()
					statePlay = true
				} else if stateLibrayMusic == "music" && statePlay == true || actualIndex > lenghtMusics {
					u.currentTime.Title = "Stop"
					u.currentTime.Label = fmt.Sprintf("00:00 / 00:00")
					ticker.Stop()
					sound.PauseSong(statePlay)
					statePlay = false
					pos = 0
					u.currentTime.Percent = 0
				}
			default:

			}

		}
		termui.Render(u.mainContainer)
	}
}
