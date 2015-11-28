package ui

import (
	"fmt"
	"time"

	"github.com/avadhutp/lazarus/events"
	"github.com/avadhutp/lazarus/geddit"

	"github.com/gizak/termui"
)

var (
	Quit     = quitWidget()
	Songs    = songsWidget()
	Title    = titleWidget()
	Download = downloadWidget()
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

func downloadWidget() *termui.List {
	d := termui.NewList()
	d.Items = []string{"Waiting to download..."}
	d.Border = false
	d.Height = 10
	d.PaddingTop, d.PaddingLeft = 1, 2

	return d
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
		t := fmt.Sprintf("[%d] %s [(%s)](fg-cyan)", i+1, s.Data.Title, s.Data.Genre)
		songs = songs[0 : i+1]
		songs[i] = t
	}

	Songs.Items = songs
	Songs.Height = len(songs) + 2

	Refresh()
}

func updateDownloader(e termui.Event) {
	lst := e.Data.(*geddit.Listing)
	i := make([]string, 2, 2)

	curr, next := getCurrAndNextSongs(lst)

	i[0] = fmt.Sprintf("Downloading now: [%s](fg-blue)", curr.Data.Title)
	i[1] = fmt.Sprintf("Next in queue: [%s](fg-magenta)", next.Data.Title)

	Download.Items = i
	Refresh()

	time.Sleep(1 * time.Second)
	DownloadSong(curr)

	if next.Data.Url != "" {
		lst.Data.Children = lst.Data.Children[1:]
		events.FireStartSongDownload(lst)
	}
}

func getCurrAndNextSongs(lst *geddit.Listing) (curr, next geddit.Children) {
	curr = lst.Data.Children[0]

	if len(lst.Data.Children) > 1 {
		next = lst.Data.Children[1]
	}

	return
}

func EventHandler() {
	termui.Handle("/sys/kbd/q", func(termui.Event) { termui.StopLoop() })

	termui.Handle(events.FinishedGedditDownload, updateSongList)
	termui.Handle(events.StartSongDownload, updateDownloader)
}
