package main

import (
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"time"
)

func main() {
	// Create terminal
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	// colors
	containerBorderColor := cell.ColorRGB24(0, 112, 255)
	containerTitleColor := cell.ColorRGB24(250, 55, 115)

	// Create container
	c, err := container.New(
		t,
		container.Border(linestyle.Round),
		container.BorderTitle("Termplay [PRESS Q TO QUIT]"),
		container.TitleColor(containerTitleColor),
		container.TitleFocusedColor(containerTitleColor),
		container.BorderTitleAlignCenter(),
		container.BorderColor(containerBorderColor),
		container.FocusedColor(containerBorderColor),
		container.SplitVertical(
			container.Left(
				container.Border(linestyle.Round),
				container.BorderTitle("Biblioteca"),
				container.MarginRightPercent(30),
			),
			container.Right(
				container.Border(linestyle.Round),
				container.BorderTitle("Musicas"),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	// Create shortcuts
	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}
	// Run
	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(1*time.Second)); err != nil {
		panic(err)
	}
}
