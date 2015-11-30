package ui

import (
	"fmt"
	"time"

	"github.com/avadhutp/lazarus/geddit"
	"github.com/gizak/termui"
)

var (
	Quit  = quitWidget()
	Songs = songsWidget()
	Title = titleWidget()
	Log   = logWidget()
)

func Refresh() { termui.Body.Align(); termui.Render(termui.Body) }

func titleWidget() *termui.Gauge {
	t := termui.NewGauge()
	t.Label = "****** Lazarus ☢ ******"
	t.Height = 1
	t.Border = false

	return t
}

func logWidget() *termui.Par {
	l := termui.NewPar("")
	l.Height = 10

	return l
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

func updateLog(msg string) {
	Log.Text = msg
	Refresh()

	time.Sleep(5 * time.Second)
}

func formatSong(el *geddit.Children) (t string) {
	s := el.Data
	status, statusFg, titleFg := " ", "fg-white", "fg-white"

	switch s.Status {
	case geddit.IsDownloading:
		status = "»"
	case geddit.Downloaded:
		status = "✔"
		statusFg, titleFg = "fg-green", "fg-green"
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

	if e.Path == FinishedRedditDownload {
		go obj.Start()
	}
}
