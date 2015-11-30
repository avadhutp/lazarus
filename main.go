package main

import (
	"os"

	"github.com/avadhutp/lazarus/geddit"
	"github.com/avadhutp/lazarus/player"
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
			termui.NewCol(6, 0, player.Log),
		),
		termui.NewRow(
			termui.NewCol(6, 0, player.Quit),
		),
	)
}

func download() {
	lst := geddit.Get()
	playerObj := player.Player{lst, []string{}}

	player.FireFinishedRedditDownload(playerObj)
}
