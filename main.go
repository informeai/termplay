package main

import (
	"fmt"
	"os"
	"path/filepath"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getLibraryOfArgs() []string {
	if len(os.Args) > 1 {
		rowsLibrary, err := getFiles(os.Args[1])
		if err != nil {
			panic(err)
		}
		return rowsLibrary
	} else {
		rowsLibrary, err := getFiles("")
		if err != nil {
			panic(err)
		}
		return rowsLibrary
	}
}

func getFiles(path string) ([]string, error) {
	var librarys []string
	err := filepath.Walk(path, func(root string, info os.FileInfo, err error) error {

		if info.IsDir() {
			librarys = append(librarys, root)
		}
		return nil

	})

	return librarys, err
}

func listLibrary(paths []string) []string {
	var rows []string
	for _, p := range paths {
		rows = append(rows, filepath.Base(p))
	}
	return rows
}

func listMusic(index int, paths []string) []string {
	var musics []string
	err := filepath.Walk(paths[index], func(root string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			musics = append(musics, filepath.Base(root))
		}
		return nil

	})
	if err != nil {
		panic(err)
	}
	return musics
}

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
	shortsCuts.TitleStyle.Fg = ui.ColorYellow
	shortsCuts.Block.BorderStyle.Fg = ui.ColorMagenta
	shortsCuts.Text = "[ Enter ](fg:green,bg:magenta)[ Select ](fg:magenta,bg:green) " +
		"[ p ](fg:green,bg:magenta)[ Play/Pause ](fg:magenta,bg:green) " +
		"[ Esc ](fg:green,bg:magenta)[ Stop ](fg:magenta,bg:green) " +
		"[ Left ](fg:green,bg:magenta)[ Library ](fg:magenta,bg:green) " +
		"[ Right ](fg:green,bg:magenta)[ Musics ](fg:magenta,bg:green) " +
		"[ + ](fg:green,bg:magenta)[ +Volume ](fg:magenta,bg:green) " +
		"[ - ](fg:green,bg:magenta)[ -Volume ](fg:magenta,bg:green) " +
		"[ q ](fg:green,bg:red)[ Exit ](fg:magenta,bg:green) "
	shortsCuts.Border = true
	//library ui
	library := widgets.NewList()
	library.Border = true
	library.Title = "Library"
	//get library
	rowsLibrary := getLibraryOfArgs()
	library.Rows = listLibrary(rowsLibrary)

	// verify row selected
	library.SelectedRowStyle = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
	library.TitleStyle.Fg = ui.ColorGreen
	library.Block.BorderStyle.Fg = ui.ColorMagenta
	// music ui
	music := widgets.NewList()
	music.Title = "Musics"
	music.Border = true
	music.TitleStyle.Fg = ui.ColorCyan
	music.Block.BorderStyle.Fg = ui.ColorMagenta
	music.Rows = listMusic(library.SelectedRow, rowsLibrary)
	//verify music

	// time
	currentTime := widgets.NewGauge()
	currentTime.Title = "current time"
	currentTime.TitleStyle.Fg = ui.ColorCyan
	currentTime.Label = "00:00/00:00"
	currentTime.LabelStyle.Fg = ui.ColorGreen
	currentTime.Percent = 50
	currentTime.BarColor = ui.ColorBlue
	currentTime.PaddingTop = 1
	currentTime.PaddingLeft = 1
	currentTime.PaddingRight = 10
	currentTime.Block.BorderStyle.Fg = ui.ColorMagenta
	// volume
	volume := widgets.NewGauge()
	volume.Title = "volume"
	volume.Percent = 40
	volume.Label = fmt.Sprint(volume.Percent, "%")
	volume.TitleStyle = ui.NewStyle(ui.ColorCyan)
	volume.Block.BorderStyle = ui.NewStyle(ui.ColorMagenta)
	volume.BarColor = ui.ColorBlue
	volume.LabelStyle = ui.NewStyle(ui.ColorGreen)
	volume.LabelStyle.Modifier = ui.Modifier(ui.AlignCenter)
	volume.PaddingTop = 1
	volume.PaddingLeft = 1
	volume.PaddingRight = 1

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
	var stateLibrayMusic = "library"
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "q":
				return
			case "j", "<Down>":
				if stateLibrayMusic == "library" {
					library.ScrollDown()
					music.Rows = listMusic(library.SelectedRow, rowsLibrary)
				} else {
					music.ScrollDown()
				}
			case "k", "<Up>":
				if stateLibrayMusic == "library" {
					library.ScrollUp()
					music.Rows = listMusic(library.SelectedRow, rowsLibrary)
				} else {
					music.ScrollUp()
				}
			case "=", "+":
				if volume.Percent >= 0 && volume.Percent < 100 {
					volume.Percent += 10
					volume.Label = fmt.Sprint(volume.Percent, "%")
				}
			case "-", "_":
				if volume.Percent > 0 && volume.Percent <= 100 {
					volume.Percent -= 10
					volume.Label = fmt.Sprint(volume.Percent, "%")
				}
			case "<Left>":
				stateLibrayMusic = "library"
				library.ScrollTop()
				library.SelectedRowStyle = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
				music.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlack)
				music.Rows = listMusic(library.SelectedRow, rowsLibrary)
				music.ScrollTop()
			case "<Right>":
				stateLibrayMusic = "music"
				library.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlack)
				music.SelectedRowStyle = ui.NewStyle(ui.ColorMagenta, ui.ColorBlack)
				music.Rows = listMusic(library.SelectedRow, rowsLibrary)
				music.ScrollTop()
			default:

			}

		}
		ui.Render(mainContainer)
	}
}
