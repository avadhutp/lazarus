package player

import (
	"fmt"

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
	t.Label = "****** Lazarus ******"
	t.Height = 1
	t.Border = false

	return t
}

func songsWidget() *termui.List {
	w := termui.NewList()
	w.Items = []string{"Downloading..."}
	w.BorderLabel = "Song list"
	w.Height = 5

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
	music := e.Data.(map[string]geddit.Children)
	songs := make([]string, 0, len(music))

	for _, s := range music {
		t := fmt.Sprintf("[ ] %s [(%s)](fg-cyan)", s.Data.Title, s.Data.Genre)
		songs = append(songs, t)
	}

	Songs.Items = songs
	Songs.Height = len(songs) + 2

	Refresh()
}
