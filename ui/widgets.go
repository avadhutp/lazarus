package ui

import (
	"fmt"

	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
)

var (
	// Quit Is a widget to display the keystroke to quit lazarus
	Quit = quitWidget()
	// Songs Is a widget to display the songlist
	Songs = songsWidget()
	// Title Is a widget to display the title of the prog
	Title = titleWidget()
)

// Refresh Re-draws the widgets after something in them changes
func Refresh() { termui.Body.Align(); termui.Render(termui.Body) }

// titleWidget Provides the title bar
func titleWidget() *termui.Gauge {
	t := termui.NewGauge()
	t.Height = 1
	t.Border = false

	return t
}

// songsWidget Provides the song list widget
func songsWidget() *termui.List {
	w := termui.NewList()
	w.Items = []string{"Downloading..."}
	w.BorderLabel = "Song list"
	w.Height = 5

	return w
}

// quitWidget Displays the key required to quit
func quitWidget() *termui.Par {
	q := termui.NewPar("Press q to quit Lazarus.")
	q.TextFgColor = termui.ColorRed
	q.Height = 1
	q.Border = false

	return q
}

// formatSong Formats the song displayed in Song widget according to its current status
func formatSong(el *geddit.Children) (t string) {
	s := el.Data
	status, statusFg, titleFg := " ", "fg-white", "fg-white"

	switch s.Status {
	case geddit.IsDownloading:
		status = "»"
	case geddit.Downloaded:
		status = "✔"
		statusFg, titleFg = "fg-green", "fg-green"
	case geddit.NotDownloaded:
		status = "✖"
		statusFg, titleFg = "fg-red", "fg-red"
	}

	t = fmt.Sprintf("|[%s](%s)|[%s](%s)", status, statusFg, s.Title, titleFg)

	return t
}

// UpdateSongList Will update the songs in the Songs widget when the corresponding event is fired.
func paintSongList(e termui.Event) {
	obj := e.Data.(Player)
	songs := make([]string, 0, len(obj.Music))

	for _, key := range obj.GetKeys() {
		s := obj.Music[key]

		t := formatSong(s)
		songs = append(songs, t)
	}

	Songs.Items = songs
	Songs.Height = len(songs) + 2

	Refresh()
}
