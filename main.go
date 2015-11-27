package main

import (
	"github.com/avadhutp/lazarus/events"
	"github.com/avadhutp/lazarus/geddit"
	"github.com/avadhutp/lazarus/ui"
	"github.com/gizak/termui"
)

func main() {
	ui.EventHandler()
	go download()

	render()
	defer termui.Close()
	ui.Refresh()
	termui.Loop()
}

func render() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, ui.Title),
		),
		termui.NewRow(
			termui.NewCol(6, 0, ui.Songs),
			termui.NewCol(6, 0, ui.Download),
		),
		termui.NewRow(
			termui.NewCol(12, 0, ui.Quit),
		),
	)
}

func download() {
	lst := geddit.Get()
	events.FireFinishedGedditDownload(&lst)
	events.FireStartSongDownload(&lst)
}
