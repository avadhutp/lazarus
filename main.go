package main

import (
	"os"

	"github.com/avadhutp/lazarus/ui"

	"github.com/avadhutp/lazarus/geddit"
	"github.com/codegangsta/cli"
	"github.com/gizak/termui"
)

const (
	name    = "lazarus"
	version = "0.0.1"
	desc    = "Lazarus: The resurrection of the revolution."
)

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = desc

	app.Action = start
	app.Run(os.Args)
}

func start(ctx *cli.Context) {
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
			termui.NewCol(6, 0, ui.Log),
		),
		termui.NewRow(
			termui.NewCol(6, 0, ui.Quit),
		),
	)
}

func download() {
	lst := geddit.Get()
	player := ui.Player{lst, []string{}}

	ui.FireFinishedRedditDownload(player)
}
