package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	// terminal dimensions
	termWidth, termHeight := ui.TerminalDimensions()
	//Create widgets
	mainContainer := ui.NewGrid()
	mainContainer.SetRect(0, 0, termWidth, termHeight)
	//short cuts ui
	shortsCuts := widgets.NewParagraph()
	shortsCuts.Title = "Keys"
	shortsCuts.Text = "[ Enter ](fg-black,bg-white)[Select](fg-black,bg-green) " +
		"[ p ](fg-black,bg-white)[Play/Pause](fg-black,bg-green) " +
		"[Esc](fg-black,bg-white)[Stop](fg-black,bg-green) " +
		"[Right](fg-black,bg-white)[+10s](fg-black,bg-green) " +
		"[Left](fg-black,bg-white)[-10s](fg-black,bg-green) " +
		"[ + ](fg-black,bg-white)[+Volume](fg-black,bg-green) " +
		"[ - ](fg-black,bg-white)[-Volume](fg-black,bg-green) " +
		"[ q ](fg-black,bg-white)[Exit](fg-black,bg-green) "
	shortsCuts.Border = true
	//library ui
	library := widgets.NewList()
	library.Border = true
	library.Title = "Library"
	// music ui
	music := widgets.NewList()
	music.Title = "Musics"
	music.Border = true
	// time
	currentTime := widgets.NewGauge()
	currentTime.Title = "current time"
	currentTime.Label = "00:00/00:00"
	// volume
	volume := widgets.NewGauge()
	volume.Title = "volume"
	volume.Label = "100%"

	// set mainContainer
	mainContainer.Set(
		ui.NewRow(1.0/9, shortsCuts),
		ui.NewRow(
			1.0/1.3,
			ui.NewCol(1.0/3, library),
			ui.NewCol(1.0/1.5, music),
		),
		ui.NewRow(
			1.0/8,
			ui.NewCol(1.0/1.5, currentTime),
			ui.NewCol(1.0/3, volume),
		),
	)

	ui.Render(mainContainer)
	// events keys
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
