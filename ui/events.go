package ui

import "github.com/gizak/termui"

const (
	FinishedRedditDownload = "/lazarus/reddit/download/done"
	songListUpdated        = "/lazarus/update/songlist"
)

func UpdatePlayer(player Player) {
	termui.SendCustomEvt(songListUpdated, player)
}

func FireFinishedRedditDownload(player Player) {
	termui.SendCustomEvt(FinishedRedditDownload, player)
}

func EventHandler() {
	termui.Handle("/sys/kbd/q", func(termui.Event) { termui.StopLoop() })

	termui.Handle(FinishedRedditDownload, paintSongList)
	termui.Handle(songListUpdated, paintSongList)
}
