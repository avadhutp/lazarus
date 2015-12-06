package ui

import "github.com/gizak/termui"

const (
	finishedRedditDownload = "/lazarus/reddit/download/done"
	songListUpdated        = "/lazarus/update/songlist"
)

// UpdatePlayer Fires the event to redraw the player
func UpdatePlayer(player Player) {
	termui.SendCustomEvt(songListUpdated, player)
}

// FireFinishedRedditDownload Fires the event to begin the playlist display/download/play process
func FireFinishedRedditDownload(player Player) {
	termui.SendCustomEvt(finishedRedditDownload, player)
}

// EventHandler Registers all the event handlers
func EventHandler() {
	termui.Handle("/sys/kbd/q", func(termui.Event) { termui.StopLoop() })

	termui.Handle(finishedRedditDownload, paintSongList)
	termui.Handle(songListUpdated, paintSongList)
}
