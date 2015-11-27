package ui

import (
	"fmt"

	"github.com/avadhutp/lazarus/events"
	"github.com/avadhutp/lazarus/geddit"

	"github.com/gizak/termui"
)

var (
	Quit  = quitWidget()
	Songs = songsWidget()
	Title = titleWidget()
)

func Refresh() { termui.Body.Align(); termui.Render(termui.Body) }

func titleWidget() *termui.Gauge {
	t := termui.NewGauge()
	t.Label = "*** Lazarus ***"
	t.Height = 1
	t.Border = false

	return t
}

func songsWidget() *termui.List {
	w := termui.NewList()
	w.Items = []string{"Downloading..."}
	w.BorderLabel = "Song list"
	w.Height = 3
	w.Y = 0

	return w
}

func quitWidget() *termui.Par {
	q := termui.NewPar("Press q to quit Lazarus.")
	q.TextFgColor = termui.ColorRed
	q.Height = 1
	q.Border = false

	return q
}

// UpdateSongList Will update the songs in the Songs widget when the corresponding event is fired.
func updateSongList(e termui.Event) {
	lst := e.Data.(*geddit.Listing)
	songs := make([]string, 0, len(lst.Data.Children))

	for i, s := range lst.Data.Children {
		t := fmt.Sprintf("[%d] %s", i, s.Data.Title)
		songs = songs[0 : i+1]
		songs[i] = t
	}

	Songs.Items = songs
	Songs.Height = len(songs)

	Refresh()
}

func EventHandler() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle(events.FinishedGedditDownload, updateSongList)
}
