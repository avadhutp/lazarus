package main

import (
	"github.com/gizak/termui"
	// "fmt"
	"github.com/avadhutp/lazarus/geddit"
	"github.com/avadhutp/lazarus/ui"
)

func main() {
	lst := geddit.Get()

	drawGrid(lst)
}

func drawGrid(lst geddit.Listing) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	songs := ui.SongsWidget(lst)
	quit := ui.QuitWidget()

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, songs),
		),
		termui.NewRow(
			termui.NewCol(6, 0, quit),
		),
	)

	termui.Body.Align()
	termui.Render(termui.Body)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}
