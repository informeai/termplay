package main

import (
	termui "github.com/gizak/termui/v3"
	ui "github.com/informeai/termplay/ui"
)

func main() {

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()
	u := ui.NewUi()
	u.Run("<path of musics>")
}
