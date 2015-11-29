package main

import (
	"github.com/avadhutp/lazarus/geddit"
	"github.com/avadhutp/lazarus/player"
	"github.com/gizak/termui"
)

func main() {
	player.EventHandler()
	go download()

	render()
	defer termui.Close()
	player.Refresh()
	termui.Loop()
}

func render() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, player.Title),
		),
		termui.NewRow(
			termui.NewCol(6, 0, player.Songs),
		),
		termui.NewRow(
			termui.NewCol(6, 0, player.Quit),
		),
	)
}

func download() {
	lst := geddit.Get()
	player.FireFinishedGedditDownload(lst)
}
