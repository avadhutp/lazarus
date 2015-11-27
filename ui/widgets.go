package ui

import "github.com/gizak/termui"

func SongsWidget() *termui.List {
	w := termui.NewList()
	w.Items = []string{"Downloading..."}
	w.BorderLabel = "Song list"
	w.Height = 1
	w.Y = 0

	return w
}

func QuitWidget() *termui.Par {
	quit := termui.NewPar("[Press q to quit Lazarus.](fg-red)")
	quit.Height = 1
	quit.Border = false

	return quit
}
