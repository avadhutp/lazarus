package ui

import (
	"fmt"

	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
)

func SongsWidget(lst geddit.Listing) *termui.List {
	songs := make([]string, 0, len(lst.Data.Children))

	for i, s := range lst.Data.Children {
		t := fmt.Sprintf("[%d] %s", i, s.Data.Title)
		songs = songs[0 : i+1]
		songs[i] = t
	}

	w := termui.NewList()
	w.Items = songs
	w.BorderLabel = "Song list"
	w.Height = len(songs)
	w.Y = 0

	return w
}

func QuitWidget() *termui.Par {
	quit := termui.NewPar("[Press q to quit Lazarus.](fg-red)")
	quit.Height = 1
	quit.Border = false

	return quit
}
