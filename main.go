package main

import (
	termui "github.com/gizak/termui/v3"
	ui "github.com/informeai/termplay/ui"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
		path = os.Getenv("PATH_MUSICS")
	}

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()
	u := ui.NewUi()
	u.Run(path)
}
